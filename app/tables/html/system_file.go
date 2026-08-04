package html

import (
	"fmt"
	"html/template"
	"strings"
)

// SystemFile 提供系统文件表格使用的 HTML、JavaScript 和 CSS 对象。
type SystemFile struct{}

// GetSystemFile 获取无状态的系统文件 HTML 对象。
func GetSystemFile() *SystemFile {
	return &SystemFile{}
}

// Status 根据状态值返回系统文件状态标签。
func (s *SystemFile) Status(value string) template.HTML {
	switch value {
	case "1":
		return template.HTML(`<span class="label label-success">正常</span>`)
	case "0":
		return template.HTML(`<span class="label label-warning">禁用</span>`)
	default:
		return template.HTML(`<span class="label label-default">删除</span>`)
	}
}

func (s *SystemFile) Preview(publicURL string, mimeType interface{}) template.HTML {
	url := template.HTMLEscapeString(publicURL)
	mime := fmt.Sprint(mimeType)
	if strings.HasPrefix(mime, "image/") {
		return template.HTML(`<a href="` + url + `" target="_blank"><img src="` + url + `" style="max-width:80px;max-height:80px;object-fit:contain;border:1px solid #eee" /></a>`)
	}
	return template.HTML(`<a href="` + url + `" target="_blank">打开</a>`)
}

func (s *SystemFile) DownloadButton(id string) template.HTML {
	id = template.HTMLEscapeString(id)
	return template.HTML(`<a class="btn btn-xs btn-primary" href="/admin/files/download/` + id + `" target="_blank">下载</a>`)
}

// QuickShareLink 渲染创建一天有效期分享链接的 URL 入口。
func (s *SystemFile) QuickShareLink(id, publicURL string) template.HTML {
	id = template.HTMLEscapeString(id)
	publicURL = template.HTMLEscapeString(publicURL)
	if publicURL == "" {
		publicURL = "创建 1 天分享链接"
	}
	return template.HTML(`<button type="button" class="btn btn-link btn-xs system-file-quick-share-link" data-file-id="` + id + `">` + publicURL + `</button>`)
}

// ShareButton 渲染打开文件分享弹窗的按钮。
func (s *SystemFile) ShareButton(id string) template.HTML {
	id = template.HTMLEscapeString(id)
	return template.HTML(`<button type="button" class="btn btn-xs btn-success system-file-share-button" data-file-id="` + id + `">分享</button>`)
}

// ShareModal 渲染分享时长选择弹窗及交互脚本。
func (s *SystemFile) ShareModal() template.HTML {
	return template.HTML(`
<div class="modal fade" id="system-file-share-modal" tabindex="-1" role="dialog" aria-hidden="true">
  <div class="modal-dialog" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
        <h4 class="modal-title">文件分享</h4>
      </div>
      <div class="modal-body">
        <p>选择有效期后将创建新的公开下载链接。</p>
        <div class="system-file-share-duration-buttons" role="group" style="display:flex;flex-wrap:wrap;gap:10px">
          <button type="button" class="btn btn-primary system-file-share-create" data-hours="24">分享 1 天</button>
          <button type="button" class="btn btn-primary system-file-share-create" data-hours="168">分享 7 天</button>
          <button type="button" class="btn btn-primary system-file-share-create" data-hours="720">分享 30 天</button>
        </div>
        <div class="system-file-share-result" style="display:none;margin-top:16px">
          <label for="system-file-share-link">分享链接</label>
          <div class="input-group">
            <input id="system-file-share-link" type="text" class="form-control system-file-share-link" readonly>
            <span class="input-group-btn">
              <button type="button" class="btn btn-default system-file-share-copy">复制链接</button>
            </span>
          </div>
        </div>
        <div class="alert alert-success system-file-share-notice" role="status" style="display:none;margin-top:12px;margin-bottom:0"></div>
      </div>
    </div>
  </div>
</div>
<script>
(function () {
  var modal = $('#system-file-share-modal');
  var noticeTimer;

  function request(path, method, body) {
    return $.ajax({
      url: path,
      type: method,
      contentType: 'application/json',
      dataType: 'json',
      data: body ? JSON.stringify(body) : undefined
    });
  }

  function errorMessage(response) {
    if (response && response.responseJSON && response.responseJSON.msg) {
      return response.responseJSON.msg;
    }
    return '请求失败，请稍后重试';
  }

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

  function resetShareResult() {
    modal.find('.system-file-share-link').val('');
    modal.find('.system-file-share-result').hide();
  }

  function showShareResult(link) {
    modal.find('.system-file-share-link').val(link);
    modal.find('.system-file-share-result').show();
  }

  function showNotice(message) {
    var notice = modal.find('.system-file-share-notice');
    window.clearTimeout(noticeTimer);
    notice.text(message).stop(true, true).fadeIn(150);
    noticeTimer = window.setTimeout(function () {
      notice.fadeOut(300);
    }, 2000);
  }

  function createShare(fileID, durationHours, button) {
    if (!fileID) {
      alert('文件标识无效');
      return;
    }
    button.prop('disabled', true);
    request('/admin/files/' + fileID + '/shares', 'POST', { durationHours: durationHours }).done(function (response) {
      if (!response || response.status === 'error' || !response.data || !response.data.downloadUrl) {
        alert(response.msg || '创建分享链接失败');
        return;
      }
      showShareResult(response.data.downloadUrl);
    }).fail(function (response) {
      alert(errorMessage(response));
    }).always(function () {
      button.prop('disabled', false);
    });
  }

  $(document).off('click.systemFileShare', '.system-file-share-button').on('click.systemFileShare', '.system-file-share-button', function () {
    modal.data('file-id', $(this).data('file-id'));
    resetShareResult();
    modal.modal('show');
  });

  $(document).off('click.systemFileShareCreate', '.system-file-share-create').on('click.systemFileShareCreate', '.system-file-share-create', function () {
    var button = $(this);
    createShare(modal.data('file-id'), Number(button.data('hours')), button);
  });

  $(document).off('click.systemFileShareCopy', '.system-file-share-copy').on('click.systemFileShareCopy', '.system-file-share-copy', function () {
    var link = modal.find('.system-file-share-link').val();
    if (!link) {
      return;
    }
    copyLink(link).then(function () {
      showNotice('分享链接已复制');
    }).catch(function () {
      window.prompt('请复制分享链接', link);
    });
  });

  $(document).off('click.systemFileQuickShare', '.system-file-quick-share-link').on('click.systemFileQuickShare', '.system-file-quick-share-link', function () {
    var button = $(this);
    modal.data('file-id', button.data('file-id'));
    resetShareResult();
    modal.modal('show');
    createShare(button.data('file-id'), 24, button);
  });
}());
</script>`)
}
