package html

import "html/template"

type TusUpload struct{}

// GetTusUploadHTML 获取TUS HTML 对象。
func GetTusUploadHTML() *TusUpload {
	return &TusUpload{}
}

// TusUploadContent 渲染单文件 Tus 上传控件。
func (t *TusUpload) TusUploadContent() template.HTML {
	return template.HTML(`
<script src="/assets/vendor/tus.min.js"></script>

<div class="system-file-tus-upload">
	<div class="input-group system-file-tus-main-row">
		<input
			type="text"
			class="form-control system-file-tus-caption"
			placeholder="请选择文件"
			readonly
		>

		<span class="input-group-btn">
			<button
				type="button"
				class="btn btn-primary system-file-tus-select"
			>
				选择文件
			</button>

			<button
				type="button"
				class="btn btn-default system-file-tus-pause"
				hidden
			>
				暂停
			</button>

			<button
				type="button"
				class="btn btn-danger system-file-tus-cancel"
				hidden
			>
				取消
			</button>
		</span>
	</div>

	<input
		type="file"
		class="system-file-tus-input"
		autocomplete="off"
 		style="display:none"
		hidden
	>

	<input
		type="hidden"
		name="tus_upload_id"
		class="system-file-tus-upload-id"
	>

	<div class="system-file-tus-progress-wrapper" hidden>
		<div class="progress system-file-tus-progress">
			<div
				class="progress-bar progress-bar-striped active"
				role="progressbar"
				aria-valuemin="0"
				aria-valuemax="100"
				aria-valuenow="0"
				style="width:0%"
			>
				0%
			</div>
		</div>
	</div>

	<div
		class="help-block system-file-tus-status"
		role="status"
		aria-live="polite"
	></div>
</div>
`)
}

// TusUploadScript 初始化单文件 Tus 客户端、暂停、取消、断点恢复和表单保护。
func (t *TusUpload) TusUploadScript(endpoint string) template.JS {
	endpoint = template.JSEscapeString(endpoint)

	return template.JS(`
setTimeout(function () {
	var roots = document.querySelectorAll(
		'.system-file-tus-upload'
	);

	roots.forEach(function (root) {
		if (root.dataset.initialized === '1') {
			return;
		}

		root.dataset.initialized = '1';

		var input = root.querySelector(
			'.system-file-tus-input'
		);

		var caption = root.querySelector(
			'.system-file-tus-caption'
		);

		var uploadID = root.querySelector(
			'.system-file-tus-upload-id'
		);

		var selectButton = root.querySelector(
			'.system-file-tus-select'
		);

		var pauseButton = root.querySelector(
			'.system-file-tus-pause'
		);

		var cancelButton = root.querySelector(
			'.system-file-tus-cancel'
		);

		var progressWrapper = root.querySelector(
			'.system-file-tus-progress-wrapper'
		);

		var progressBar = root.querySelector(
			'.system-file-tus-progress .progress-bar'
		);

		var status = root.querySelector(
			'.system-file-tus-status'
		);

		var form = root.closest('form');

		var recordID = form
			? form.querySelector('[name="id"]')
			: null;

		var hasExistingFile = !!(
			recordID &&
			recordID.value &&
			recordID.value !== '0'
		);

		var upload = null;
		var currentFile = null;

		// empty、uploading、paused、success、error
		var state = 'empty';

		function setStatus(message, type) {
			status.textContent = message || '';

			status.classList.toggle(
				'text-danger',
				type === 'error'
			);

			status.classList.toggle(
				'text-success',
				type === 'success'
			);

			status.classList.toggle(
				'text-warning',
				type === 'warning'
			);
		}

		function setSubmitEnabled(enabled) {
			if (!form) {
				return;
			}

			form.querySelectorAll(
				'button[type="submit"], input[type="submit"]'
			).forEach(function (button) {
				button.disabled = !enabled;
			});
		}

		function updateOriginalName(fileName) {
			if (!form) {
				return;
			}

			var originalName = form.querySelector(
				'[name="original_name"]'
			);

			if (originalName) {
				originalName.value = fileName || '';
			}
		}

		function setProgress(percentage) {
			var value = Number(percentage);

			if (!Number.isFinite(value)) {
				value = 0;
			}

			value = Math.max(
				0,
				Math.min(100, value)
			);

			var text = value.toFixed(1) + '%';

			if (value === 0 || value === 100) {
				text = value.toFixed(0) + '%';
			}

			progressBar.style.width = value + '%';
			progressBar.textContent = text;

			progressBar.setAttribute(
				'aria-valuenow',
				String(value)
			);
		}

		function showUploadingControls() {
			pauseButton.hidden = false;
			cancelButton.hidden = false;
			progressWrapper.hidden = false;
		}

		function hideUploadingControls() {
			pauseButton.hidden = true;
			cancelButton.hidden = true;
			progressWrapper.hidden = true;
		}

		function setState(nextState, message) {
			state = nextState;

			root.classList.toggle(
				'is-uploading',
				nextState === 'uploading'
			);

			root.classList.toggle(
				'is-paused',
				nextState === 'paused'
			);

			root.classList.toggle(
				'is-success',
				nextState === 'success'
			);

			root.classList.toggle(
				'is-error',
				nextState === 'error'
			);

			switch (nextState) {
			case 'uploading':
				showUploadingControls();

				pauseButton.textContent = '暂停';
				pauseButton.disabled = false;
				cancelButton.disabled = false;
				selectButton.disabled = true;

				progressBar.classList.add(
					'progress-bar-striped',
					'active'
				);

				setStatus(
					message || '正在上传文件',
					''
				);
				break;

			case 'paused':
				showUploadingControls();

				pauseButton.textContent = '继续';
				pauseButton.disabled = false;
				cancelButton.disabled = false;
				selectButton.disabled = false;

				progressBar.classList.remove('active');

				setStatus(
					message || '上传已暂停',
					'warning'
				);
				break;

			case 'error':
				showUploadingControls();

				pauseButton.textContent = '继续';
				pauseButton.disabled = false;
				cancelButton.disabled = false;
				selectButton.disabled = false;

				progressBar.classList.remove('active');

				setStatus(
					message || '上传失败',
					'error'
				);
				break;

			case 'success':
				hideUploadingControls();

				pauseButton.textContent = '暂停';
				pauseButton.disabled = true;
				cancelButton.disabled = true;
				selectButton.disabled = false;

				progressBar.classList.remove(
					'progress-bar-striped',
					'active'
				);

				setStatus(
					message || '上传完成，可以保存表单',
					'success'
				);
				break;

			default:
				hideUploadingControls();

				pauseButton.textContent = '暂停';
				pauseButton.disabled = true;
				cancelButton.disabled = true;
				selectButton.disabled = false;

				progressBar.classList.add(
					'progress-bar-striped'
				);

				progressBar.classList.remove('active');
			}
		}

		function completedUploadID(uploadURL) {
			if (!uploadURL) {
				return '';
			}

			var parsed = new URL(
				uploadURL,
				window.location.href
			);

			var parts = parsed.pathname
				.replace(/\/$/, '')
				.split('/');

			return decodeURIComponent(
				parts[parts.length - 1] || ''
			);
		}

		function resetUploadState(keepExistingFile) {
			upload = null;
			currentFile = null;
			state = 'empty';

			uploadID.value = '';
			input.value = '';
			caption.value = '';

			setProgress(0);
			setState('empty');

			if (keepExistingFile && hasExistingFile) {
				caption.placeholder = '已保留原文件';

				setStatus(
					'未选择新文件，将保留原文件',
					'success'
				);

				setSubmitEnabled(true);
				return;
			}

			caption.placeholder = '请选择文件';
			updateOriginalName('');

			setSubmitEnabled(false);
		}

		function terminateCurrentUpload(terminateRemote) {
			if (!upload) {
				return Promise.resolve();
			}

			var currentUpload = upload;
			upload = null;

			return currentUpload
				.abort(terminateRemote === true)
				.catch(function () {
					// 服务端终止失败时，仍允许前端恢复状态。
				});
		}

		function createUpload(file) {
			if (
				!window.tus ||
				!window.tus.Upload
			) {
				setState(
					'error',
					'tus 客户端加载失败，请刷新页面重试'
				);

				setSubmitEnabled(hasExistingFile);
				return null;
			}

			return new window.tus.Upload(file, {
				endpoint: '` + endpoint + `',

				retryDelays: [
					0,
					1000,
					3000,
					5000,
					10000
				],

				metadata: {
					filename: file.name,
					filetype: file.type || ''
				},

				removeFingerprintOnSuccess: false,

				onError: function (error) {
					if (!currentFile) {
						return;
					}

					setState(
						'error',
						'上传中断：' + error.message
					);

					setSubmitEnabled(hasExistingFile);
				},

				onProgress: function (
					uploadedBytes,
					totalBytes
				) {
					if (!currentFile) {
						return;
					}

					var percentage = totalBytes > 0
						? uploadedBytes / totalBytes * 100
						: 0;

					setProgress(percentage);

					if (state !== 'uploading') {
						setState(
							'uploading',
							'正在上传 ' +
								percentage.toFixed(1) +
								'%'
						);
					} else {
						setStatus(
							'正在上传 ' +
								percentage.toFixed(1) +
								'%',
							''
						);
					}
				},

				onSuccess: function () {
					if (!currentFile || !upload) {
						return;
					}

					var id = completedUploadID(
						upload.url || ''
					);

					if (!id) {
						setState(
							'error',
							'上传完成，但未获取到上传标识'
						);

						setSubmitEnabled(hasExistingFile);
						return;
					}

					uploadID.value = id;

					caption.value = currentFile.name;
					updateOriginalName(currentFile.name);

					setProgress(100);

					// 上传完成后隐藏进度条、暂停和取消按钮。
					setState(
						'success',
						'上传完成，可以保存表单'
					);

					setSubmitEnabled(true);
				}
			});
		}

		function startUpload(file) {
			currentFile = file;
			uploadID.value = '';

			caption.value = file.name;
			caption.title = file.name;

			setProgress(0);
			setSubmitEnabled(false);

			upload = createUpload(file);

			if (!upload) {
				return;
			}

			setState(
				'uploading',
				'正在准备上传'
			);

			upload
				.findPreviousUploads()
				.then(function (previousUploads) {
					if (
						previousUploads &&
						previousUploads.length > 0
					) {
						upload.resumeFromPreviousUpload(
							previousUploads[0]
						);
					}

					upload.start();
				})
				.catch(function () {
					upload.start();
				});
		}

		function selectFile(file) {
			if (!file) {
				return;
			}

			// 单文件模式：选择新文件时终止旧上传。
			terminateCurrentUpload(true).then(function () {
				startUpload(file);
			});
		}

		selectButton.addEventListener(
			'click',
			function () {
				if (state === 'uploading') {
					return;
				}

				// 确保可以重复选择同一个文件。
				input.value = '';
				input.click();
			}
		);

		input.addEventListener(
			'change',
			function () {
				var file = input.files &&
					input.files.length > 0
					? input.files[0]
					: null;

				if (!file) {
					return;
				}

				selectFile(file);
			}
		);

		pauseButton.addEventListener(
			'click',
			function () {
				if (!upload || !currentFile) {
					return;
				}

				if (state === 'uploading') {
					pauseButton.disabled = true;

					upload.abort()
						.then(function () {
							setState(
								'paused',
								'上传已暂停，可以继续上传'
							);

							setSubmitEnabled(
								hasExistingFile
							);
						})
						.catch(function (error) {
							pauseButton.disabled = false;

							setStatus(
								'暂停上传失败：' +
									error.message,
								'error'
							);
						});

					return;
				}

				if (
					state === 'paused' ||
					state === 'error'
				) {
					setSubmitEnabled(false);

					setState(
						'uploading',
						'正在继续上传'
					);

					upload.start();
				}
			}
		);

		cancelButton.addEventListener(
			'click',
			function () {
				cancelButton.disabled = true;
				pauseButton.disabled = true;

				terminateCurrentUpload(true).then(function () {
					resetUploadState(hasExistingFile);
				});
			}
		);

		if (form) {
			form.addEventListener(
				'submit',
				function (event) {
					var hasNewUpload =
						!!uploadID.value;

					if (state === 'uploading') {
						event.preventDefault();
						event.stopImmediatePropagation();

						setStatus(
							'请等待文件上传完成后再保存',
							'error'
						);

						return;
					}

					if (state === 'paused') {
						event.preventDefault();
						event.stopImmediatePropagation();

						setStatus(
							'当前上传已暂停，请继续上传或取消',
							'error'
						);

						return;
					}

					if (state === 'error') {
						event.preventDefault();
						event.stopImmediatePropagation();

						setStatus(
							'文件上传失败，请继续上传或取消',
							'error'
						);

						return;
					}

					if (
						!hasExistingFile &&
						!hasNewUpload
					) {
						event.preventDefault();
						event.stopImmediatePropagation();

						setStatus(
							'请先选择并完成文件上传',
							'error'
						);
					}
				},
				true
			);
		}

		resetUploadState(hasExistingFile);
	});
}, 0);`)
}

// TusUploadStyle 返回单文件 Tus 上传控件样式。
func (t *TusUpload) TusUploadStyle() template.CSS {
	return template.CSS(`
.system-file-tus-upload {
	display: block;
	width: 100%;
}

.system-file-tus-main-row {
	display: table;
	width: 100%;
}

.system-file-tus-caption {
	height: 36px;
	background-color: #fff;
	color: #606266;
	cursor: default;
}

.system-file-tus-caption[readonly] {
	background-color: #fff;
}

.system-file-tus-main-row .input-group-btn {
	width: 1%;
	white-space: nowrap;
}

.system-file-tus-main-row .btn {
	height: 36px;
	padding-right: 16px;
	padding-left: 16px;
}

.system-file-tus-main-row .btn[hidden] {
	display: none !important;
}

.system-file-tus-select {
	border-color: #409eff;
	background-color: #409eff;
	color: #fff;
}

.system-file-tus-select:hover,
.system-file-tus-select:focus {
	border-color: #66b1ff;
	background-color: #66b1ff;
	color: #fff;
}

.system-file-tus-pause {
	color: #606266;
}

.system-file-tus-cancel {
	color: #fff;
}

.system-file-tus-progress-wrapper {
	width: 100%;
	margin-top: 10px;
}

.system-file-tus-progress-wrapper[hidden] {
	display: none !important;
}

.system-file-tus-progress {
	height: 8px;
	margin: 0;
	overflow: hidden;
	border-radius: 100px;
	background-color: #ebeef5;
	box-shadow: none;
}

.system-file-tus-progress .progress-bar {
	min-width: 0;
	height: 8px;
	border-radius: 100px;
	background-color: #409eff;
	font-size: 0;
	line-height: 8px;
	box-shadow: none;
	transition: width 0.2s ease;
}

.system-file-tus-upload.is-paused
.system-file-tus-progress .progress-bar {
	background-color: #e6a23c;
}

.system-file-tus-upload.is-error
.system-file-tus-progress .progress-bar {
	background-color: #f56c6c;
}

.system-file-tus-status {
	min-height: 20px;
	margin: 5px 0 0;
	font-size: 12px;
	line-height: 20px;
}

.system-file-tus-status.text-success {
	color: #67c23a;
}

.system-file-tus-status.text-warning {
	color: #e6a23c;
}

.system-file-tus-status.text-danger {
	color: #f56c6c;
}

@media (max-width: 576px) {
	.system-file-tus-main-row {
		display: flex;
		flex-wrap: wrap;
		gap: 8px;
	}

	.system-file-tus-main-row .form-control {
		width: 100%;
		flex: 1 0 100%;
		border-radius: 4px;
	}

	.system-file-tus-main-row .input-group-btn {
		display: flex;
		width: 100%;
		gap: 8px;
	}

	.system-file-tus-main-row .input-group-btn .btn {
		flex: 1;
		border-radius: 4px;
	}
}
`)
}

// UploadHelp 返回系统文件上传字段的帮助文案。
func (t *TusUpload) UploadHelp() template.HTML {
	return template.HTML("支持暂停和断点续传；新增时必传，编辑时不选择新文件则保留原文件")
}
