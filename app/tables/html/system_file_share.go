package html

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

// SystemFileShare 提供系统文件分享表格使用的 HTML 对象。
type SystemFileShare struct{}

// GetSystemFileShare 获取无状态的系统文件分享 HTML 对象。
func GetSystemFileShare() *SystemFileShare {
	return &SystemFileShare{}
}

// CopyButton 渲染复制分享下载链接的按钮。
func (s *SystemFileShare) CopyButton(token string) template.HTML {
	token = template.HTMLEscapeString(token)
	return template.HTML(`<button type="button" class="btn btn-xs btn-primary system-file-share-copy-button" data-download-path="/shares/` + token + `/download">复制链接</button>`)
}

// CopyScript 返回分享链接复制按钮的页面交互脚本。
func (s *SystemFileShare) CopyScript() template.HTML {
	return template.HTML(`
<script>
(function () {
  function copyLink(link) {
    if (navigator.clipboard && window.isSecureContext) {
      return navigator.clipboard.writeText(link);
    }
    var input = document.createElement('textarea');
    input.value = link;
    input.style.position = 'fixed';
    input.style.opacity = '0';
    document.body.appendChild(input);
    input.select();
    var copied = document.execCommand('copy');
    document.body.removeChild(input);
    return copied ? Promise.resolve() : Promise.reject(new Error('copy failed'));
  }

  $(document).off('click.systemFileShareTableCopy', '.system-file-share-copy-button').on('click.systemFileShareTableCopy', '.system-file-share-copy-button', function () {
    var button = $(this);
    var path = button.data('download-path');
    if (!path) {
      return;
    }
    var link = new URL(path, window.location.origin).toString();
    var text = button.text();
    button.prop('disabled', true);
    copyLink(link).then(function () {
      button.text('已复制');
      window.setTimeout(function () {
        button.text(text);
      }, 1500);
      button.prop('disabled', false);
    }, function () {
      window.prompt('请复制下载链接', link);
      button.prop('disabled', false);
    });
  });
}());
</script>`)
}

// Status 根据撤销和到期状态渲染状态标签。
func (s *SystemFileShare) Status(row map[string]interface{}) template.HTML {
	if row["revoked_at"] != nil && fmt.Sprint(row["revoked_at"]) != "" {
		return template.HTML(`<span class="label label-default">已撤销</span>`)
	}
	expiresAt, err := parseShareTime(fmt.Sprint(row["expires_at"]))
	if err != nil || !expiresAt.After(time.Now()) {
		return template.HTML(`<span class="label label-warning">已过期</span>`)
	}
	return template.HTML(`<span class="label label-success">有效</span>`)
}

func parseShareTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	for _, layout := range []string{time.DateTime, time.RFC3339, time.RFC3339Nano} {
		if result, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return result, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid share expiration time")
}
