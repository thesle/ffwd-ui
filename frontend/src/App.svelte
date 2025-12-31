<script>
  import { onMount } from 'svelte';
  import { EventsOn } from '../wailsjs/runtime/runtime';
  import * as App from '../wailsjs/go/main/App';

  let inputFile = '';
  let outputFile = '';
  let fileInfo = null;
  let operation = 'trim_start';
  
  // Time inputs (H:M:S.s)
  let trimStartH = 0, trimStartM = 0, trimStartS = 0;
  let trimLengthH = 0, trimLengthM = 1, trimLengthS = 0;
  let trimRangeStartH = 0, trimRangeStartM = 0, trimRangeStartS = 0;
  let trimRangeEndH = 0, trimRangeEndM = 1, trimRangeEndS = 0;
  let paddingStartH = 0, paddingStartM = 0, paddingStartS = 0;
  let paddingEndH = 0, paddingEndM = 0, paddingEndS = 0;
  
  // Audio
  let audioFormat = 'mp3';
  
  // Format conversion
  let targetFormat = 'mp4';
  
  // Resolution
  let resolutionPreset = '1920x1080';
  let customWidth = 1920;
  let customHeight = 1080;
  
  // Volume
  let volumePercent = 100;
  
  // Crop
  let cropWidth = 1280;
  let cropHeight = 720;
  let cropX = 0;
  let cropY = 0;
  
  // Bitrate
  let videoBitrate = '2M';
  let audioBitrate = '192k';
  let useTwoPass = false;
  
  // Hardware acceleration
  let useHardwareAccel = false;
  let hardwareEncoder = 'none';
  
  // UI state
  let thumbnail = null;
  let darkMode = false;
  
  let diskSpace = [];
  let commandPreview = '';
  let isRunning = false;
  let progress = 0;
  let progressMessage = '';
  let errorMessage = '';
  let successMessage = '';
  let isDragging = false;

  onMount(async () => {
    loadDiskSpace();
    
    // Load dark mode preference
    const savedDarkMode = localStorage.getItem('darkMode');
    if (savedDarkMode) {
      darkMode = savedDarkMode === 'true';
    }
    
    // Detect hardware encoder
    try {
      const detected = await App.DetectHardwareEncoder();
      if (detected && detected !== 'none') {
        hardwareEncoder = detected;
      }
    } catch (err) {
      console.log('Could not detect hardware encoder:', err);
    }

    // Keyboard shortcuts
    const handleKeyboard = (e) => {
      if (e.key === ' ' && !isRunning && inputFile && outputFile) {
        e.preventDefault();
        execute();
      } else if (e.key === 'Escape' && isRunning) {
        e.preventDefault();
        cancel();
      } else if ((e.ctrlKey || e.metaKey) && e.key === 'o') {
        e.preventDefault();
        selectInputFile();
      }
    };
    
    window.addEventListener('keydown', handleKeyboard);

    EventsOn('ffmpeg:progress', (data) => {
      progress = data.percent || 0;
      progressMessage = data.message || '';
    });

    EventsOn('ffmpeg:complete', () => {
      isRunning = false;
      progress = 100;
      successMessage = 'Operation completed successfully!';
      setTimeout(() => {
        successMessage = '';
        progress = 0;
        progressMessage = '';
      }, 3000);
    });

    EventsOn('ffmpeg:error', (error) => {
      isRunning = false;
      errorMessage = error;
      progress = 0;
      progressMessage = '';
      setTimeout(() => errorMessage = '', 5000);
    });
    
    return () => {
      window.removeEventListener('keydown', handleKeyboard);
    };
  });

  async function loadDiskSpace() {
    try {
      diskSpace = await App.GetDiskSpace();
    } catch (err) {
      console.error('Failed to load disk space:', err);
    }
  }

  function handleDragOver(e) {
    e.preventDefault();
    e.stopPropagation();
    isDragging = true;
  }

  function handleDragLeave(e) {
    e.preventDefault();
    e.stopPropagation();
    isDragging = false;
  }

  async function handleDrop(e) {
    e.preventDefault();
    e.stopPropagation();
    isDragging = false;

    const files = e.dataTransfer?.files;
    if (files && files.length > 0) {
      const file = files[0];
      const filePath = file.path;
      
      if (filePath && filePath.startsWith('/')) {
        try {
          inputFile = filePath;
          fileInfo = await App.GetFileInfo(filePath);
          outputFile = await App.GetDefaultOutputName(filePath, operation);
          updateCommandPreview();
        } catch (err) {
          errorMessage = 'Failed to process dropped file: ' + err;
        }
      } else {
        errorMessage = 'Drag-and-drop not supported in this environment. Please click to browse for files.';
        setTimeout(() => errorMessage = '', 5000);
      }
    }
  }

  async function selectInputFile() {
    try {
      const file = await App.SelectInputFile();
      if (file) {
        inputFile = file;
        fileInfo = await App.GetFileInfo(file);
        outputFile = await App.GetDefaultOutputName(file, operation);
        
        // Load thumbnail for video files
        if (fileInfo && fileInfo.width > 0) {
          try {
            thumbnail = await App.ExtractThumbnail(file);
          } catch (err) {
            console.error('Could not extract thumbnail:', err);
            thumbnail = null;
            // Don't show error to user - thumbnail is optional
          }
        } else {
          thumbnail = null;
        }
        
        updateCommandPreview();
      }
    } catch (err) {
      errorMessage = 'Failed to select file: ' + err;
    }
  }

  async function selectOutputFile() {
    try {
      const file = await App.SelectOutputFile(outputFile);
      if (file) {
        outputFile = file;
        updateCommandPreview();
      }
    } catch (err) {
      errorMessage = 'Failed to select output file: ' + err;
    }
  }

  function timeToSeconds(h, m, s) {
    return h * 3600 + m * 60 + parseFloat(s);
  }

  async function updateCommandPreview() {
    if (!inputFile || !outputFile) {
      commandPreview = '';
      return;
    }

    try {
      const params = {};
      
      switch(operation) {
        case 'trim_start':
          params.seconds = timeToSeconds(trimStartH, trimStartM, trimStartS);
          break;
        case 'trim_length':
          params.duration = timeToSeconds(trimLengthH, trimLengthM, trimLengthS);
          break;
        case 'trim_range':
          params.start_seconds = timeToSeconds(trimRangeStartH, trimRangeStartM, trimRangeStartS);
          params.end_seconds = timeToSeconds(trimRangeEndH, trimRangeEndM, trimRangeEndS);
          break;
        case 'extract_audio':
          params.format = audioFormat;
          break;
        case 'convert_format':
          break;
        case 'change_resolution':
          if (resolutionPreset === 'custom') {
            params.width = customWidth;
            params.height = customHeight;
          } else {
            const [w, h] = resolutionPreset.split('x').map(Number);
            params.width = w;
            params.height = h;
          }
          if (useHardwareAccel) {
            params.hw_accel = hardwareEncoder;
          }
          break;
        case 'adjust_volume':
          params.volume_percent = volumePercent;
          break;
        case 'crop_video':
          params.width = cropWidth;
          params.height = cropHeight;
          params.x = cropX;
          params.y = cropY;
          break;
        case 'adjust_bitrate':
          params.video_bitrate = videoBitrate;
          params.audio_bitrate = audioBitrate;
          if (useHardwareAccel) {
            params.hw_accel = hardwareEncoder;
          }
          params.two_pass = useTwoPass;
          break;
        case 'add_padding':
          params.start_seconds = timeToSeconds(paddingStartH, paddingStartM, paddingStartS);
          params.end_seconds = timeToSeconds(paddingEndH, paddingEndM, paddingEndS);
          break;
      }

      commandPreview = await App.PreviewCommand(operation, inputFile, outputFile, params);
    } catch (err) {
      console.error('Failed to preview command:', err);
    }
  }

  async function execute() {
    if (!inputFile || !outputFile) {
      errorMessage = 'Please select input and output files';
      return;
    }

    errorMessage = '';
    successMessage = '';
    isRunning = true;
    progress = 0;

    try {
      switch(operation) {
        case 'trim_start':
          await App.TrimStart(inputFile, outputFile, timeToSeconds(trimStartH, trimStartM, trimStartS));
          break;
        case 'trim_length':
          await App.TrimToLength(inputFile, outputFile, timeToSeconds(trimLengthH, trimLengthM, trimLengthS));
          break;
        case 'trim_range':
          await App.TrimRange(inputFile, outputFile, 
            timeToSeconds(trimRangeStartH, trimRangeStartM, trimRangeStartS),
            timeToSeconds(trimRangeEndH, trimRangeEndM, trimRangeEndS));
          break;
        case 'extract_audio':
          await App.ExtractAudio(inputFile, outputFile, audioFormat);
          break;
        case 'convert_format':
          await App.ConvertFormat(inputFile, outputFile);
          break;
        case 'change_resolution':
          const width = resolutionPreset === 'custom' ? customWidth : parseInt(resolutionPreset.split('x')[0]);
          const height = resolutionPreset === 'custom' ? customHeight : parseInt(resolutionPreset.split('x')[1]);
          await App.ChangeResolution(inputFile, outputFile, width, height, useHardwareAccel ? hardwareEncoder : 'none');
          break;
        case 'adjust_volume':
          await App.AdjustVolume(inputFile, outputFile, volumePercent);
          break;
        case 'crop_video':
          await App.CropVideo(inputFile, outputFile, cropWidth, cropHeight, cropX, cropY);
          break;
        case 'adjust_bitrate':
          await App.AdjustBitrate(inputFile, outputFile, videoBitrate, audioBitrate, useHardwareAccel ? hardwareEncoder : 'none', useTwoPass);
          break;
        case 'add_padding':
          await App.AddPadding(inputFile, outputFile,
            timeToSeconds(paddingStartH, paddingStartM, paddingStartS),
            timeToSeconds(paddingEndH, paddingEndM, paddingEndS));
          break;
      }
    } catch (err) {
      errorMessage = 'Operation failed: ' + err;
      isRunning = false;
    }
  }

  async function cancel() {
    try {
      await App.CancelOperation();
      isRunning = false;
      progress = 0;
      progressMessage = 'Operation cancelled';
    } catch (err) {
      errorMessage = 'Failed to cancel: ' + err;
    }
  }

  function formatBytes(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
  }

  function formatDuration(seconds) {
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = Math.floor(seconds % 60);
    return `${h}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
  }

  function toggleDarkMode() {
    darkMode = !darkMode;
    localStorage.setItem('darkMode', darkMode.toString());
  }
  
  function updateOutputExtension(newExt) {
    if (!outputFile) return;
    const lastDot = outputFile.lastIndexOf('.');
    if (lastDot > 0) {
      outputFile = outputFile.substring(0, lastDot) + '.' + newExt;
    }
  }

  // Auto-update output extension when format changes (must run before command preview)
  $: if (operation === 'convert_format' && targetFormat && outputFile) {
    updateOutputExtension(targetFormat);
  }
  
  $: if (operation === 'extract_audio' && audioFormat && outputFile) {
    updateOutputExtension(audioFormat);
  }

  $: {
    operation;
    trimStartH; trimStartM; trimStartS;
    trimLengthH; trimLengthM; trimLengthS;
    trimRangeStartH; trimRangeStartM; trimRangeStartS;
    trimRangeEndH; trimRangeEndM; trimRangeEndS;
    paddingStartH; paddingStartM; paddingStartS;
    paddingEndH; paddingEndM; paddingEndS;
    audioFormat; targetFormat; resolutionPreset; customWidth; customHeight;
    volumePercent; cropWidth; cropHeight; cropX; cropY;
    videoBitrate; audioBitrate; useTwoPass; useHardwareAccel; hardwareEncoder;
    outputFile; // Make sure this triggers after outputFile updates
    updateCommandPreview();
  }

  $: if (inputFile && operation) {
    App.GetDefaultOutputName(inputFile, operation).then(name => outputFile = name);
  }
</script>

<div class="app-container" class:dark-mode={darkMode} style="height: 100vh; display: flex; flex-direction: column;">
  <section class="app-section" style="flex: 1; overflow-y: auto; padding-bottom: 100px;">
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem;">
      <h1 class="title" style="margin: 0;">FFwd UI</h1>
      <button class="button is-small" on:click={toggleDarkMode} title="Toggle dark mode">
        {#if darkMode}
          ☀️
        {:else}
          🌙
        {/if}
      </button>
    </div>

    {#if errorMessage}
      <div class="notification is-danger is-light">
        <button class="delete" on:click={() => errorMessage = ''}></button>
        {errorMessage}
      </div>
    {/if}

    {#if successMessage}
      <div class="notification is-success is-light">
        <button class="delete" on:click={() => successMessage = ''}></button>
        {successMessage}
      </div>
    {/if}

    <div class="columns">
      <!-- Left Column: Active Controls -->
      <div class="column is-8">
        <div class="box">
      <div class="field">
        <label class="label">Input File</label>
        <div class="drop-zone" class:active={isDragging} on:click={selectInputFile} on:dragover={handleDragOver} on:dragleave={handleDragLeave} on:drop={handleDrop}>
          {#if inputFile}
            <p class="has-text-weight-bold">{inputFile}</p>
            {#if fileInfo}
              <p class="is-size-7 has-text-grey">
                Size: {formatBytes(fileInfo.size)} | Duration: {formatDuration(fileInfo.duration)} | Format: {fileInfo.format}
              </p>
            {/if}
          {:else}
            <p>Click to select a file</p>
          {/if}
        </div>
      </div>

      <div class="field">
        <label class="label">Operation</label>
        <div class="control">
          <label class="radio"><input type="radio" bind:group={operation} value="trim_start"> Trim Start</label>
          <label class="radio"><input type="radio" bind:group={operation} value="trim_length"> Trim to Length</label>
          <label class="radio"><input type="radio" bind:group={operation} value="trim_range"> Trim Range</label>
          <label class="radio"><input type="radio" bind:group={operation} value="extract_audio"> Extract Audio</label>
          <label class="radio"><input type="radio" bind:group={operation} value="convert_format"> Convert Format</label>
          <label class="radio"><input type="radio" bind:group={operation} value="change_resolution"> Change Resolution</label>
          <label class="radio"><input type="radio" bind:group={operation} value="adjust_volume"> Adjust Volume</label>
          <label class="radio"><input type="radio" bind:group={operation} value="crop_video"> Crop Video</label>
          <label class="radio"><input type="radio" bind:group={operation} value="adjust_bitrate"> Adjust Bitrate</label>
          <label class="radio"><input type="radio" bind:group={operation} value="add_padding"> Add Padding</label>
        </div>
      </div>

      {#if operation === 'trim_start'}
        <div class="field">
          <label class="label">Time to Remove from Start</label>
          <div class="field is-grouped">
            <div class="control"><input class="input" type="number" bind:value={trimStartH} min="0" placeholder="H" style="width: 70px;"></div>
            <div class="control"><span class="is-size-4">:</span></div>
            <div class="control"><input class="input" type="number" bind:value={trimStartM} min="0" max="59" placeholder="M" style="width: 70px;"></div>
            <div class="control"><span class="is-size-4">:</span></div>
            <div class="control"><input class="input" type="number" bind:value={trimStartS} min="0" max="59.99" step="0.01" placeholder="S.ss" style="width: 80px;"></div>
          </div>
        </div>
      {/if}

      {#if operation === 'trim_length'}
        <div class="field">
          <label class="label">Target Duration</label>
          <div class="field is-grouped">
            <div class="control"><input class="input" type="number" bind:value={trimLengthH} min="0" placeholder="H" style="width: 70px;"></div>
            <div class="control"><span class="is-size-4">:</span></div>
            <div class="control"><input class="input" type="number" bind:value={trimLengthM} min="0" max="59" placeholder="M" style="width: 70px;"></div>
            <div class="control"><span class="is-size-4">:</span></div>
            <div class="control"><input class="input" type="number" bind:value={trimLengthS} min="0" max="59.99" step="0.01" placeholder="S.ss" style="width: 80px;"></div>
          </div>
        </div>
      {/if}

      {#if operation === 'trim_range'}
        <div class="field">
          <label class="label">Start Time</label>
          <div class="field is-grouped">
            <div class="control"><input class="input" type="number" bind:value={trimRangeStartH} min="0" placeholder="H" style="width: 70px;"></div>
            <div class="control"><span class="is-size-4">:</span></div>
            <div class="control"><input class="input" type="number" bind:value={trimRangeStartM} min="0" max="59" placeholder="M" style="width: 70px;"></div>
            <div class="control"><span class="is-size-4">:</span></div>
            <div class="control"><input class="input" type="number" bind:value={trimRangeStartS} min="0" max="59.99" step="0.01" placeholder="S.ss" style="width: 80px;"></div>
          </div>
        </div>
        <div class="field">
          <label class="label">End Time</label>
          <div class="field is-grouped">
            <div class="control"><input class="input" type="number" bind:value={trimRangeEndH} min="0" placeholder="H" style="width: 70px;"></div>
            <div class="control"><span class="is-size-4">:</span></div>
            <div class="control"><input class="input" type="number" bind:value={trimRangeEndM} min="0" max="59" placeholder="M" style="width: 70px;"></div>
            <div class="control"><span class="is-size-4">:</span></div>
            <div class="control"><input class="input" type="number" bind:value={trimRangeEndS} min="0" max="59.99" step="0.01" placeholder="S.ss" style="width: 80px;"></div>
          </div>
        </div>
      {/if}

      {#if operation === 'extract_audio'}
        <div class="field">
          <label class="label">Audio Format</label>
          <div class="control">
            <div class="select">
              <select bind:value={audioFormat}>
                <option value="mp3">MP3</option>
                <option value="aac">AAC</option>
                <option value="wav">WAV</option>
                <option value="flac">FLAC</option>
              </select>
            </div>
          </div>
        </div>
      {/if}

      {#if operation === 'convert_format'}
        <div class="field">
          <label class="label">Target Format</label>
          <div class="control">
            <div class="select">
              <select bind:value={targetFormat}>
                <option value="mp4">MP4</option>
                <option value="mkv">MKV</option>
                <option value="avi">AVI</option>
                <option value="mov">MOV</option>
                <option value="webm">WebM</option>
              </select>
            </div>
          </div>
          <p class="help">Change output file extension to match target format</p>
        </div>
      {/if}

      {#if operation === 'change_resolution'}
        <div class="field">
          <label class="label">Resolution</label>
          <div class="control">
            <div class="select">
              <select bind:value={resolutionPreset}>
                <option value="3840x2160">4K (3840x2160)</option>
                <option value="1920x1080">1080p (1920x1080)</option>
                <option value="1280x720">720p (1280x720)</option>
                <option value="854x480">480p (854x480)</option>
                <option value="640x360">360p (640x360)</option>
                <option value="custom">Custom</option>
              </select>
            </div>
          </div>
        </div>
        {#if resolutionPreset === 'custom'}
          <div class="field is-grouped">
            <div class="control is-expanded">
              <input class="input" type="number" bind:value={customWidth} min="1" placeholder="Width">
            </div>
            <div class="control is-expanded">
              <input class="input" type="number" bind:value={customHeight} min="1" placeholder="Height">
            </div>
          </div>
        {/if}
      {/if}

      {#if operation === 'adjust_volume'}
        <div class="field">
          <label class="label">Volume ({volumePercent}%)</label>
          <div class="control">
            <input class="input" type="range" bind:value={volumePercent} min="0" max="200" step="5">
          </div>
          <p class="help">50% = half volume, 100% = original, 200% = double volume</p>
        </div>
      {/if}

      {#if operation === 'crop_video'}
        <div class="field">
          <label class="label">Crop Dimensions</label>
          <div class="field is-grouped">
            <div class="control is-expanded">
              <input class="input" type="number" bind:value={cropWidth} min="1" placeholder="Width">
            </div>
            <div class="control is-expanded">
              <input class="input" type="number" bind:value={cropHeight} min="1" placeholder="Height">
            </div>
          </div>
        </div>
        <div class="field">
          <label class="label">Crop Position (X, Y)</label>
          <div class="field is-grouped">
            <div class="control is-expanded">
              <input class="input" type="number" bind:value={cropX} min="0" placeholder="X offset">
            </div>
            <div class="control is-expanded">
              <input class="input" type="number" bind:value={cropY} min="0" placeholder="Y offset">
            </div>
          </div>
        </div>
      {/if}

      {#if operation === 'add_padding'}
        <div class="field">
          <label class="label">Padding at Start</label>
          <div class="field is-grouped">
            <div class="control"><input class="input" type="number" bind:value={paddingStartH} min="0" placeholder="H" style="width: 70px;"></div>
            <div class="control"><span class="is-size-4">:</span></div>
            <div class="control"><input class="input" type="number" bind:value={paddingStartM} min="0" max="59" placeholder="M" style="width: 70px;"></div>
            <div class="control"><span class="is-size-4">:</span></div>
            <div class="control"><input class="input" type="number" bind:value={paddingStartS} min="0" max="59.99" step="0.01" placeholder="S.ss" style="width: 80px;"></div>
          </div>
        </div>
        <div class="field">
          <label class="label">Padding at End</label>
          <div class="field is-grouped">
            <div class="control"><input class="input" type="number" bind:value={paddingEndH} min="0" placeholder="H" style="width: 70px;"></div>
            <div class="control"><span class="is-size-4">:</span></div>
            <div class="control"><input class="input" type="number" bind:value={paddingEndM} min="0" max="59" placeholder="M" style="width: 70px;"></div>
            <div class="control"><span class="is-size-4">:</span></div>
            <div class="control"><input class="input" type="number" bind:value={paddingEndS} min="0" max="59.99" step="0.01" placeholder="S.ss" style="width: 80px;"></div>
          </div>
        </div>
      {/if}

      {#if operation === 'adjust_bitrate'}
        <div class="field">
          <label class="label">Video Bitrate</label>
          <div class="control">
            <input class="input" type="text" bind:value={videoBitrate} placeholder="e.g., 2M, 5000k">
          </div>
          <p class="help">Leave empty to copy without re-encoding</p>
        </div>
        <div class="field">
          <label class="label">Audio Bitrate</label>
          <div class="control">
            <input class="input" type="text" bind:value={audioBitrate} placeholder="e.g., 192k, 320k">
          </div>
          <p class="help">Leave empty to copy without re-encoding</p>
        </div>
        <div class="field">
          <label class="checkbox">
            <input type="checkbox" bind:checked={useTwoPass}>
            Use two-pass encoding (better quality)
          </label>
        </div>
      {/if}
      
      {#if (operation === 'change_resolution' || operation === 'adjust_bitrate') && hardwareEncoder !== 'none'}
        <div class="field">
          <label class="checkbox">
            <input type="checkbox" bind:checked={useHardwareAccel}>
            Use hardware acceleration ({hardwareEncoder})
          </label>
          <p class="help">Faster encoding using GPU</p>
        </div>
      {/if}

      <div class="field">
        <label class="label">Output File</label>
        <div class="field has-addons">
          <div class="control is-expanded">
            <input class="input" type="text" bind:value={outputFile} placeholder="Output file path">
          </div>
          <div class="control">
            <button class="button is-info" on:click={selectOutputFile}>Browse</button>
          </div>
        </div>
      </div>

      <div class="field is-grouped">
        <div class="control">
          <button class="button is-primary" on:click={execute} disabled={isRunning || !inputFile || !outputFile}>
            Execute
          </button>
        </div>
        <div class="control">
          <button class="button is-danger" on:click={cancel} disabled={!isRunning}>
            Cancel
          </button>
        </div>
      </div>

      {#if isRunning || progress > 0}
        <div class="field">
          <label class="label">Progress</label>
          <progress class="progress is-primary" value={progress} max="100">{progress}%</progress>
          {#if progressMessage}
            <p class="is-size-7 has-text-grey">{progressMessage}</p>
          {/if}
        </div>
      {/if}
        </div>
      </div>

      <!-- Right Column: Passive Display -->
      <div class="column is-4">
        <div class="box info-panel">
          <h2 class="subtitle">System Information</h2>
          
          {#if fileInfo}
            <div class="field">
              <label class="label">File Details</label>
              {#if thumbnail}
                <figure class="image" style="margin-bottom: 1rem;">
                  <img src="{thumbnail}" alt="Video thumbnail" style="border-radius: 4px; max-width: 100%;" />
                </figure>
              {/if}
              <div class="content">
                <p class="is-size-7"><strong>Size:</strong> {formatBytes(fileInfo.size)}</p>
                <p class="is-size-7"><strong>Duration:</strong> {formatDuration(fileInfo.duration)}</p>
                <p class="is-size-7"><strong>Format:</strong> {fileInfo.format}</p>
                {#if fileInfo.codec}
                  <p class="is-size-7"><strong>Codec:</strong> {fileInfo.codec}</p>
                {/if}
                {#if fileInfo.width && fileInfo.height}
                  <p class="is-size-7"><strong>Resolution:</strong> {fileInfo.width}x{fileInfo.height}</p>
                {/if}
              </div>
            </div>
            <hr>
          {/if}

          {#if diskSpace.length > 0 && diskSpace.length <= 10}
            <div class="field">
              <label class="label">Available Disk Space</label>
              <div class="content">
                {#each diskSpace as mount}
                  <div class="disk-space-item">
                    <p class="is-size-7 has-text-weight-bold">{mount.path}</p>
                    <p class="is-size-7 has-text-grey">
                      {formatBytes(mount.available)} free of {formatBytes(mount.total)}
                    </p>
                    <progress class="progress is-small is-info" value={mount.available} max={mount.total}></progress>
                  </div>
                {/each}
              </div>
            </div>
          {/if}
        </div>
      </div>
    </div>
  </section>

  <div class="command-preview">
    <strong>Command:</strong> {commandPreview || 'Select files and operation to preview command'}
  </div>
</div>

<style>
  .app-container {
    max-width: 100%;
    margin: 0;
    padding: 0;
  }

  .app-section {
    padding: 1.5rem 1rem;
  }

  .drop-zone {
    min-height: 80px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
  }

  .radio {
    margin-right: 1rem;
  }

  .info-panel {
    background-color: #f8f9fa;
    border-left: 4px solid #3273dc;
    position: sticky;
    top: 1rem;
  }

  .disk-space-item {
    margin-bottom: 1rem;
  }

  .disk-space-item:last-child {
    margin-bottom: 0;
  }

  .columns {
    margin-top: 1rem;
  }
  
  /* Dark Mode Styles */
  .dark-mode {
    background-color: #1a1a1a;
    color: #e0e0e0;
  }
  
  .dark-mode .title,
  .dark-mode .subtitle,
  .dark-mode .label {
    color: #e0e0e0;
  }
  
  .dark-mode .box {
    background-color: #2d2d2d;
    color: #e0e0e0;
  }
  
  .dark-mode .info-panel {
    background-color: #252525;
    border-left-color: #5a9fd4;
  }
  
  .dark-mode .input,
  .dark-mode .select select {
    background-color: #3a3a3a;
    color: #e0e0e0;
    border-color: #555;
  }
  
  .dark-mode .input::placeholder {
    color: #888;
  }
  
  .dark-mode .drop-zone {
    background-color: #2d2d2d;
    border-color: #555;
  }
  
  .dark-mode .notification.is-danger {
    background-color: #4a2626;
    color: #ffb3b3;
  }
  
  .dark-mode .notification.is-success {
    background-color: #264a26;
    color: #b3ffb3;
  }
  
  .dark-mode .command-preview {
    background-color: #2d2d2d;
    color: #e0e0e0;
    border-top-color: #555;
  }
  
  .dark-mode .has-text-grey {
    color: #aaa !important;
  }
  
  .dark-mode .progress::-webkit-progress-value {
    background-color: #5a9fd4;
  }
  
  .dark-mode .progress::-moz-progress-bar {
    background-color: #5a9fd4;
  }
</style>
