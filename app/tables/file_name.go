package tables

import (
	"html/template"
	"path"
	"strings"

	"github.com/GoAdminGroup/go-admin/template/types"
)

// ConfigureFileNameDisplay 启用文件名列的固定宽度布局、提示和复制交互。
func ConfigureFileNameDisplay(info *types.InfoPanel) {
	info.SetTableFixed().SetFooterHtml(fileNameDisplayFooter())
}

// DisplayFileName 将文件名渲染为主体省略、后缀保留的列表内容。
func DisplayFileName(value types.FieldModel) interface{} {
	return renderFileName(value.Value)
}

// splitFileName 分离文件名主体和扩展名，点文件保持原始名称。
func splitFileName(fileName string) (string, string) {
	extension := path.Ext(fileName)
	name := strings.TrimSuffix(fileName, extension)
	if name == "" {
		return fileName, ""
	}
	return name, extension
}

// renderFileName 渲染带完整文件名提示和复制操作的文件名内容。
func renderFileName(fileName string) template.HTML {
	name, extension := splitFileName(fileName)
	escapedFileName := template.HTMLEscapeString(fileName)
	escapedName := template.HTMLEscapeString(name)
	escapedExtension := template.HTMLEscapeString(extension)

	return template.HTML(`<div style="display:flex;align-items:center;width:100%;min-width:0;gap:4px">` +
		`<span class="table-file-name-tooltip" tabindex="0" data-toggle="tooltip" data-placement="top" title="` + escapedFileName + `" style="display:flex;flex:1 1 auto;min-width:0">` +
		`<span style="flex:1 1 auto;min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap">` + escapedName + `</span>` +
		`<span style="flex:0 0 auto;white-space:nowrap">` + escapedExtension + `</span>` +
		`</span>` +
		`<button type="button" class="btn btn-link btn-xs table-file-name-copy" data-file-name="` + escapedFileName + `" title="复制文件名" style="padding:0"><i class="fa fa-copy"></i></button>` +
		`</div>`)
}

// fileNameDisplayFooter 返回文件名提示和复制操作的样式与脚本。
func fileNameDisplayFooter() template.HTML {
	return template.HTML(`
<style>
.table-file-name-tooltip-popup .tooltip-inner {
  max-width: 480px;
  word-break: break-all;
}
</style>
<script>
(function () {
  function copyFileName(fileName) {
    if (navigator.clipboard && window.isSecureContext) {
      return navigator.clipboard.writeText(fileName);
    }
    var input = document.createElement('textarea');
    input.value = fileName;
    input.style.position = 'fixed';
    input.style.opacity = '0';
    document.body.appendChild(input);
    input.select();
    var copied = document.execCommand('copy');
    document.body.removeChild(input);
    return copied ? Promise.resolve() : Promise.reject(new Error('copy failed'));
  }

  $('body').tooltip({
    selector: '.table-file-name-tooltip',
    container: 'body',
    placement: 'top',
    html: false,
    template: '<div class="tooltip table-file-name-tooltip-popup" role="tooltip"><div class="tooltip-arrow"></div><div class="tooltip-inner"></div></div>'
  });

  $(document).off('click.tableFileNameCopy', '.table-file-name-copy').on('click.tableFileNameCopy', '.table-file-name-copy', function () {
    var button = $(this);
    var fileName = button.data('file-name');
    if (!fileName) {
      return;
    }
    copyFileName(fileName).then(function () {
      button.tooltip({title: '已复制', placement: 'bottom'}).tooltip('show');
    }).catch(function () {
      window.prompt('请复制文件名', fileName);
    });
  });
}());
</script>`)
}
