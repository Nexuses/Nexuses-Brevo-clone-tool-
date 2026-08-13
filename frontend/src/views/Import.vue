<template>
  <section class="import import-brevo bv-page" :class="{ 'is-file-wizard': step === 'file' }">
    <crm-subnav v-if="step !== 'file'" />

    <header v-if="isFree() && step !== 'file'" class="import-brevo__header">
      <div class="import-brevo__title-row">
        <button
          v-if="step !== 'chooser'"
          type="button"
          class="import-brevo__back"
          aria-label="Back"
          @click="step = 'chooser'"
        >
          <svg width="18" height="18" viewBox="0 0 18 18" fill="none" aria-hidden="true">
            <path d="M11.5 4.5L7 9l4.5 4.5" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
        </button>
        <div>
          <h1 class="import-brevo__title">
            Import contacts for bulk creation or updating
          </h1>
          <p class="import-brevo__lead">
            Create, update, or blocklist contacts in bulk. Keep in mind you must have your contacts'
            consent to send them campaigns.
          </p>
        </div>
      </div>
    </header>
    <b-loading :active="isLoading" />

    <!-- Brevo-style method chooser -->
    <section v-if="isFree() && step === 'chooser'" class="import-brevo__chooser">
      <div class="import-brevo__methods">
        <button type="button" class="import-brevo__method" @click="openStep('file')">
          <span class="import-brevo__method-icon" aria-hidden="true">
            <svg width="28" height="28" viewBox="0 0 28 28" fill="none">
              <rect x="6" y="4" width="14" height="18" rx="2" stroke="#0B8F4D" stroke-width="1.6" />
              <path d="M10 12h8M10 16h6" stroke="#0B8F4D" stroke-width="1.6" stroke-linecap="round" />
              <path d="M18 20v4M16 22h4" stroke="#0B8F4D" stroke-width="1.6" stroke-linecap="round" />
            </svg>
          </span>
          <span class="import-brevo__method-title">Import from a file</span>
          <span class="import-brevo__method-desc">Import your contacts from a csv, xlsx, or txt file.</span>
        </button>

        <button type="button" class="import-brevo__method" @click="openStep('paste')">
          <span class="import-brevo__method-icon" aria-hidden="true">
            <svg width="28" height="28" viewBox="0 0 28 28" fill="none">
              <rect x="8" y="8" width="12" height="14" rx="2" stroke="#0B8F4D" stroke-width="1.6" />
              <rect x="5" y="5" width="12" height="14" rx="2" stroke="#0B8F4D" stroke-width="1.6" fill="#fff" />
            </svg>
          </span>
          <span class="import-brevo__method-title">Copy-paste</span>
          <span class="import-brevo__method-desc">Paste the contacts as text from a spreadsheet or a similar list.</span>
        </button>
      </div>
    </section>

    <section v-if="isFree() && step === 'file'" class="import-file">
      <header class="import-file__head">
        <div class="import-file__head-row">
          <h1 class="import-file__title">Import contacts from a file</h1>
          <button type="button" class="import-file__close" aria-label="Close" @click="closeFileWizard">
            <svg width="18" height="18" viewBox="0 0 18 18" fill="none" aria-hidden="true">
              <path d="M4.5 4.5l9 9M13.5 4.5l-9 9" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" />
            </svg>
          </button>
        </div>
        <p class="import-file__lead">
          Upload a file containing all your contacts and their information.
          This is particularly useful when you have a large number of contacts to import.
        </p>
        <a
          class="import-file__learn"
          href="https://listmonk.app/docs/apis/import.html"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn more on how to import your contacts
          <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
            <path d="M4.5 2.5H2.75A1.25 1.25 0 0 0 1.5 3.75v5.5A1.25 1.25 0 0 0 2.75 10.5h5.5A1.25 1.25 0 0 0 9.5 9.25V7.5" stroke="currentColor" stroke-width="1.3" stroke-linecap="round" />
            <path d="M6.5 1.5H10.5V5.5M10.5 1.5L5.5 6.5" stroke="currentColor" stroke-width="1.3" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
        </a>
      </header>

      <div class="import-file__step import-file__step--first is-active">
        <h2 class="import-file__step-title">
          <span class="import-file__num">1</span>
          Upload your file
        </h2>
        <div class="import-file__step-meta">
          <p>Select a file containing your contacts to import.</p>
          <button type="button" class="import-file__example" @click="downloadExample">
            <svg width="14" height="14" viewBox="0 0 14 14" fill="none" aria-hidden="true">
              <path d="M7 2v7M4.5 6.5L7 9l2.5-2.5M2.5 11.5h9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
            Download example file (.csv)
          </button>
        </div>

        <input
          id="import-file-input"
          ref="fileInput"
          class="is-sr-only"
          type="file"
          accept=".csv,.txt,.zip,text/csv,text/plain,application/zip"
          aria-label="Upload contacts file"
          @change="onFileInput"
        />

        <button
          type="button"
          class="import-file__drop"
          :class="{ 'is-dragover': fileDragOver, 'has-file': !!form.file }"
          @click="openFilePicker"
          @dragover.prevent="fileDragOver = true"
          @dragleave.prevent="fileDragOver = false"
          @drop.prevent="onFileDrop"
        >
          <span class="import-file__cloud" aria-hidden="true">
            <svg width="88" height="64" viewBox="0 0 88 64" fill="none">
              <path
                fill="#C5D3E6"
                d="M70.8 26.4C69.2 16.6 60.6 9.2 50.4 9.2c-3.4 0-6.6.8-9.4 2.3C37.4 5.4
                  30.8 2 23.4 2 11.6 2 2 11.6 2 23.4c0 1.8.2 3.5.7 5.1C1.1 30.8 0 34.2 0 38c0
                  8.8 7.2 16 16 16h55.2C80.4 54 88 46.4 88 37.2c0-6.8-4.1-12.7-10-15.3.5-1.8.8-3.6.8-5.5z"
              />
              <path
                d="M44 20v20M44 20l-8.5 8.5M44 20l8.5 8.5"
                stroke="#fff"
                stroke-width="3.2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </span>
          <template v-if="form.file">
            <span class="import-file__drop-title">{{ form.file.name }}</span>
            <span class="import-file__drop-sub">Click to replace this file</span>
          </template>
          <template v-else>
            <span class="import-file__drop-title">Select your file or drag and drop it here</span>
            <span class="import-file__drop-sub">.csv, .xlsx or .txt</span>
          </template>
        </button>

        <p v-if="form.file" class="import-file__file-actions">
          <button type="button" class="import-file__clear" @click.stop="clearFile">Remove file</button>
        </p>

        <p class="import-file__privacy">
          <svg width="14" height="14" viewBox="0 0 14 14" fill="none" aria-hidden="true">
            <rect x="2.5" y="6" width="9" height="6.5" rx="1.4" stroke="currentColor" stroke-width="1.3" />
            <path d="M4.5 6V4.4a2.5 2.5 0 0 1 5 0V6" stroke="currentColor" stroke-width="1.3" stroke-linecap="round" />
          </svg>
          We don't sell, rent or use your database for any commercial purposes.
        </p>
      </div>

      <div class="import-file__step" :class="{ 'is-active': fileStep >= 2 }">
        <h2 class="import-file__step-title">
          <span class="import-file__num">2</span>
          Mapping data
        </h2>
        <template v-if="fileStep >= 2">
          <p v-if="fileParseError" class="import-file__error">{{ fileParseError }}</p>
          <template v-else-if="isZipFile">
            <p>ZIP archive selected. Columns inside the CSV will be mapped automatically.</p>
            <button type="button" class="import-file__continue" @click="confirmMapping">Continue</button>
          </template>
          <template v-else>
            <p>
              {{ fileRowCount }} contact{{ fileRowCount === 1 ? '' : 's' }} detected.
              Match each column in your file to a contact field.
            </p>
            <div class="import-file__map">
              <div class="import-file__map-head">
                <span>Column in your file</span>
                <span>Import as</span>
              </div>
              <div v-for="(col, idx) in fileColumns" :key="idx" class="import-file__map-row">
                <span class="import-file__map-name">{{ col.original }}</span>
                <select
                  v-model="col.mapped"
                  class="import-file__map-select"
                  :aria-label="'Import column ' + col.original + ' as'"
                >
                  <option value="">Do not import</option>
                  <option value="email">Email</option>
                  <option value="firstname">First name</option>
                  <option value="lastname">Last name</option>
                  <option value="name">Full name</option>
                  <option value="company">Company</option>
                  <option value="attributes">Attributes (JSON)</option>
                </select>
              </div>
            </div>
            <button type="button" class="import-file__continue" @click="confirmMapping">Continue</button>
          </template>
        </template>
      </div>

      <div class="import-file__step" :class="{ 'is-active': fileStep >= 3 }">
        <h2 class="import-file__step-title">
          <span class="import-file__num">3</span>
          Select a list
        </h2>
        <template v-if="fileStep >= 3">
          <div class="import-file__options">
            <div>
              <p class="import-file__opt-label">{{ $t('import.mode') }}</p>
              <b-radio v-model="form.mode" name="file-mode" native-value="subscribe" data-cy="check-subscribe">
                {{ $t('import.subscribe') }}
              </b-radio>
              <b-radio v-model="form.mode" name="file-mode" native-value="blocklist" data-cy="check-blocklist">
                {{ $t('import.blocklist') }}
              </b-radio>
            </div>
            <div>
              <p class="import-file__opt-label">{{ $t('globals.fields.status') }}</p>
              <template v-if="form.mode === 'subscribe'">
                <b-radio v-model="form.subStatus" name="file-subStatus" native-value="unconfirmed" data-cy="check-unconfirmed">
                  {{ $t('subscribers.status.unconfirmed') }}
                </b-radio>
                <b-radio v-model="form.subStatus" name="file-subStatus" native-value="confirmed" data-cy="check-confirmed">
                  {{ $t('subscribers.status.confirmed') }}
                </b-radio>
              </template>
              <b-radio v-else v-model="form.subStatus" name="file-subStatus" native-value="unsubscribed" data-cy="check-unsubscribed">
                {{ $t('subscribers.status.unsubscribed') }}
              </b-radio>
            </div>
          </div>

          <list-selector
            v-if="form.mode === 'subscribe'"
            :label="$t('globals.terms.lists')"
            :placeholder="$t('import.listSubHelp')"
            :message="$t('import.listSubHelp')"
            v-model="form.lists"
            :selected="form.lists"
            :all="lists.results"
          />

          <div v-if="form.mode === 'subscribe'" class="import-file__toggles">
            <b-field :label="$t('import.overwriteUserInfo')" :message="$t('import.overwriteUserInfoHelp')">
              <b-switch v-model="form.overwriteUserInfo" name="overwriteUserInfo" data-cy="overwrite-user-info" />
            </b-field>
            <b-field :label="$t('import.overwriteSubStatus')" :message="$t('import.overwriteSubStatusHelp')">
              <b-switch v-model="form.overwriteSubStatus" name="overwriteSubStatus" data-cy="overwrite-sub-status" />
            </b-field>
          </div>

          <div class="import-file__actions">
            <b-button
              type="is-dark"
              class="btn-new"
              :disabled="!canSubmitFile"
              :loading="isProcessing"
              @click="onUpload"
            >
              {{ $t('import.upload') }}
            </b-button>
          </div>
        </template>
      </div>
    </section><!-- upload //-->

    <section v-if="isFree() && step === 'paste'" class="import-brevo__paste">
      <div class="import-brevo__paste-head">
        <h2>Import contacts with copy/paste</h2>
        <p>
          Copy your contacts and their information from a file and paste them here.
          This is useful when you have a small number of contacts to import.
        </p>
      </div>

      <div class="import-brevo__paste-card">
        <h3>1) Import your data</h3>
        <p class="import-brevo__paste-sub">Copy and paste your contacts and their information from a file.</p>
        <label v-if="!pasteHeaders.length" class="is-sr-only" for="paste-import-box">Paste contacts data</label>
        <textarea
          v-if="!pasteHeaders.length"
          id="paste-import-box"
          v-model="form.paste"
          class="import-brevo__paste-box"
          rows="10"
          placeholder="CONTACT ID,EMAIL,FIRSTNAME,LASTNAME&#10;12345,emma@example.com,Emma,Dubois"
          @paste="onPasteData"
          @input="onPasteInput"
        />
        <div v-else class="import-brevo__paste-preview">
          <div class="import-brevo__paste-preview-meta">
            <span>{{ pasteRows.length }} row(s) · {{ pasteHeaders.length }} column(s)</span>
            <button type="button" class="import-brevo__paste-reset" @click="clearPasteGrid">
              Paste different data
            </button>
          </div>
          <div class="import-brevo__paste-preview-scroll">
            <table>
              <caption class="is-sr-only">Pasted contacts preview</caption>
              <thead>
                <tr>
                  <th v-for="(h, hi) in pasteHeaders" :key="'ph-' + hi" :title="h">{{ h }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(row, ri) in pastePreviewRows" :key="'pr-' + ri">
                  <td
                    v-for="(cell, ci) in paddedPasteRow(row)"
                    :key="'pc-' + ri + '-' + ci"
                    :title="cell"
                  >
                    {{ cell }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <p v-if="pasteRows.length > 25" class="import-brevo__paste-preview-note">
            Showing the first 25 rows. All {{ pasteRows.length }} rows will be imported.
          </p>
        </div>
        <p class="import-brevo__paste-note">
          We don't sell, rent or use your database for any commercial purposes.
        </p>
        <button
          type="button"
          class="import-brevo__check-btn"
          :disabled="!hasPasteData"
          @click="checkPasteData"
        >
          Check the data
        </button>
      </div>

      <div class="import-brevo__paste-step" :class="{ 'is-active': pasteChecked }">
        <h3>2) Mapping data</h3>
        <p v-if="pasteChecked">{{ pasteSummary }}</p>
        <p v-else>Upload your data first.</p>
      </div>

      <div class="import-brevo__paste-step" :class="{ 'is-active': pasteChecked }">
        <h3>3) Select a list</h3>
        <div v-if="pasteChecked">
          <list-selector
            v-if="form.mode === 'subscribe'"
            :label="$t('globals.terms.lists')"
            :placeholder="$t('import.listSubHelp')"
            :message="$t('import.listSubHelp')"
            v-model="form.lists"
            :selected="form.lists"
            :all="lists.results"
          />
          <div class="import-brevo__paste-actions">
            <b-button
              type="is-dark"
              class="btn-new"
              :disabled="!canSubmitPaste"
              :loading="isProcessing"
              @click="onUpload"
            >
              {{ $t('import.upload') }}
            </b-button>
          </div>
        </div>
        <p v-else>Check your data to continue.</p>
      </div>
    </section>

    <section v-if="isRunning() || isDone()" class="wrap status box has-text-centered">
      <b-progress :value="progress" show-value type="is-success" />
      <br />
      <p
        :class="['is-size-5', 'is-capitalized', { 'has-text-success': status.status === 'finished' }, { 'has-text-danger': (status.status === 'failed' || status.status === 'stopped') }]">
        {{ status.status }}
      </p>

      <p>{{ $t('import.recordsCount', { num: status.imported, total: status.total }) }}</p>
      <br />

      <p>
        <b-button @click="stopImport" :loading="isProcessing" icon-left="file-upload-outline" type="is-dark" class="btn-new">
          {{ isDone() ? $t('import.importDone') : $t('import.stopImport') }}
        </b-button>
      </p>
      <br />

      <div class="import-logs">
        <log-view :lines="logs" :loading="false" />
      </div>
    </section>
  </section>
</template>

<script>
import Vue from 'vue';
import { mapState } from 'vuex';
import ListSelector from '../components/ListSelector.vue';
import LogView from '../components/LogView.vue';
import CrmSubnav from '../components/CrmSubnav.vue';

export default Vue.extend({
  components: {
    ListSelector,
    LogView,
    CrmSubnav,
  },

  props: {
    data: { type: Object, default: () => { } },
    isEditing: { type: Boolean, default: false },
  },

  data() {
    return {
      step: 'chooser',
      form: {
        mode: 'subscribe',
        subStatus: 'unconfirmed',
        delim: ',',
        lists: [],
        overwriteUserInfo: false,
        overwriteSubStatus: false,
        file: null,
        paste: '',
      },
      example: '',
      pasteChecked: false,
      pasteSummary: '',
      pasteHeaders: [],
      pasteRows: [],
      fileStep: 1,
      fileDragOver: false,
      fileColumns: [],
      fileRowCount: 0,
      fileRawText: '',
      fileParseError: '',
      isZipFile: false,

      isLoading: true,
      isProcessing: false,
      status: { status: '' },
      logs: [],
      pollID: null,
    };
  },

  watch: {
    'form.mode': function formMode() {
      this.$nextTick(() => {
        if (this.form.mode === 'subscribe') {
          this.form.subStatus = 'unconfirmed';
        } else {
          this.form.subStatus = 'unsubscribed';
        }
      });
    },
  },

  methods: {
    openStep(nextStep) {
      this.step = nextStep;
      if (nextStep === 'paste') {
        this.clearPasteGrid();
      }
      if (nextStep === 'file') {
        this.resetFileWizard();
      }
    },

    closeFileWizard() {
      this.resetFileWizard();
      this.step = 'chooser';
    },

    resetFileWizard() {
      this.fileStep = 1;
      this.fileDragOver = false;
      this.fileColumns = [];
      this.fileRowCount = 0;
      this.fileRawText = '';
      this.fileParseError = '';
      this.isZipFile = false;
      this.form.file = null;
      if (this.$refs.fileInput) {
        this.$refs.fileInput.value = '';
      }
    },

    openFilePicker() {
      if (this.$refs.fileInput) {
        this.$refs.fileInput.click();
      }
    },

    onFileInput(e) {
      const file = e.target.files && e.target.files[0];
      if (file) this.setImportFile(file);
    },

    onFileDrop(e) {
      this.fileDragOver = false;
      const file = e.dataTransfer && e.dataTransfer.files && e.dataTransfer.files[0];
      if (file) this.setImportFile(file);
    },

    setImportFile(file) {
      const name = (file.name || '').toLowerCase();
      if (name.endsWith('.xlsx') || name.endsWith('.xls')) {
        this.$utils.toast('Excel files are not supported. Please export as .csv or .txt and try again.', 'is-danger');
        return;
      }
      const ok = name.endsWith('.csv') || name.endsWith('.txt') || name.endsWith('.zip');
      if (!ok) {
        this.$utils.toast('Please choose a .csv, .txt, or .zip file.', 'is-danger');
        return;
      }
      this.form.file = file;
      this.isZipFile = name.endsWith('.zip');
      this.fileParseError = '';
      this.fileColumns = [];
      this.fileRowCount = 0;
      this.fileRawText = '';
      this.fileStep = 2;

      if (this.isZipFile) return;

      const reader = new FileReader();
      reader.onload = () => {
        const text = typeof reader.result === 'string' ? reader.result : '';
        this.parseImportFile(text);
      };
      reader.onerror = () => {
        this.fileParseError = 'Could not read this file. Try another CSV.';
      };
      reader.readAsText(file);
    },

    parseCsvLine(line, delim) {
      const out = [];
      let cur = '';
      let inQuotes = false;
      for (let i = 0; i < line.length; i += 1) {
        const ch = line[i];
        if (ch === '"') {
          if (inQuotes && line[i + 1] === '"') {
            cur += '"';
            i += 1;
          } else {
            inQuotes = !inQuotes;
          }
        } else if (ch === delim && !inQuotes) {
          out.push(cur.trim());
          cur = '';
        } else {
          cur += ch;
        }
      }
      out.push(cur.trim());
      return out;
    },

    detectDelim(line) {
      const counts = {
        '\t': (line.match(/\t/g) || []).length,
        ';': (line.match(/;/g) || []).length,
        ',': (line.match(/,/g) || []).length,
      };
      let best = ',';
      let max = -1;
      Object.keys(counts).forEach((d) => {
        if (counts[d] > max) {
          max = counts[d];
          best = d;
        }
      });
      if (max <= 0) return ',';
      return best;
    },

    clearPasteGrid() {
      this.pasteHeaders = [];
      this.pasteRows = [];
      this.pasteChecked = false;
      this.pasteSummary = '';
      this.form.paste = '';
      this.form.delim = ',';
    },

    paddedPasteRow(row) {
      const out = Array.isArray(row) ? row.slice() : [];
      while (out.length < this.pasteHeaders.length) {
        out.push('');
      }
      return out.slice(0, this.pasteHeaders.length);
    },

    attribKeyFromHeader(header) {
      return String(header || '')
        .toLowerCase()
        .replace(/[^a-z0-9]+/g, '_')
        .replace(/^_|_$/g, '');
    },

    csvEscape(value) {
      const text = String(value == null ? '' : value);
      if (/[",\n\r]/.test(text)) {
        return `"${text.replace(/"/g, '""')}"`;
      }
      return text;
    },

    parsePastedHtmlTable(html) {
      if (!html || html.toLowerCase().indexOf('<table') === -1) return null;
      const doc = new DOMParser().parseFromString(html, 'text/html');
      const table = doc.querySelector('table');
      if (!table) return null;
      const rows = Array.from(table.querySelectorAll('tr')).map((tr) => (
        Array.from(tr.querySelectorAll('th,td')).map((cell) => (
          String(cell.innerText || '').replace(/\s+/g, ' ').trim()
        ))
      )).filter((row) => row.some((cell) => cell));
      if (rows.length < 2) return null;
      return { headers: rows[0], rows: rows.slice(1), delim: ',' };
    },

    parsePastedText(text) {
      const raw = String(text || '')
        .replace(/^\uFEFF/, '')
        .replace(/\r\n/g, '\n')
        .replace(/\r/g, '\n')
        .trim();
      if (!raw) return null;
      const lines = raw.split(/\r?\n/).filter((line) => line.trim().length > 0);
      if (lines.length < 2) return null;
      const delim = this.detectDelim(lines[0]);
      const headers = this.parseCsvLine(lines[0], delim);
      if (headers.filter((h) => h).length < 2) return null;
      const rows = lines.slice(1).map((line) => this.parseCsvLine(line, delim));
      return { headers, rows, delim };
    },

    applyPasteGrid(parsed) {
      if (!parsed || !parsed.headers || parsed.headers.length < 2) return false;
      this.pasteHeaders = parsed.headers.map((h) => String(h || '').trim());
      this.pasteRows = parsed.rows || [];
      this.form.delim = parsed.delim || ',';
      this.form.paste = [this.pasteHeaders, ...this.pasteRows]
        .map((row) => row.map((cell) => this.csvEscape(cell)).join(','))
        .join('\n');
      this.pasteChecked = false;
      this.pasteSummary = '';
      return true;
    },

    onPasteData(e) {
      const clip = e.clipboardData;
      if (!clip) return;
      const html = clip.getData('text/html');
      const text = clip.getData('text/plain');
      const parsed = this.parsePastedHtmlTable(html) || this.parsePastedText(text);
      if (!parsed) return;
      e.preventDefault();
      this.applyPasteGrid(parsed);
    },

    onPasteInput() {
      if (!(this.form.paste || '').trim()) {
        this.pasteHeaders = [];
        this.pasteRows = [];
        this.pasteChecked = false;
        this.pasteSummary = '';
        return;
      }
      if ((this.form.paste.match(/\t/g) || []).length < 2) return;
      const parsed = this.parsePastedText(this.form.paste);
      if (parsed) this.applyPasteGrid(parsed);
    },

    buildImportCsvFromPaste() {
      if (!this.pasteHeaders.length || !this.pasteRows.length) return '';
      const mapped = this.pasteHeaders.map((h) => this.guessMappedField(h));
      const knownIdx = {
        email: -1,
        firstname: -1,
        lastname: -1,
        name: -1,
        company: -1,
      };
      mapped.forEach((key, i) => {
        if (key && knownIdx[key] === -1) knownIdx[key] = i;
      });
      if (knownIdx.email === -1) return '';

      const extraIdx = mapped.map((key, i) => (key ? -1 : i)).filter((i) => i >= 0);
      const outHeaders = ['email'];
      if (knownIdx.firstname !== -1) outHeaders.push('firstname');
      if (knownIdx.lastname !== -1) outHeaders.push('lastname');
      if (knownIdx.name !== -1) outHeaders.push('name');
      if (knownIdx.company !== -1) outHeaders.push('company');
      if (extraIdx.length) outHeaders.push('attributes');

      const lines = [outHeaders.join(',')];
      this.pasteRows.forEach((row) => {
        const cells = outHeaders.map((key) => {
          if (key !== 'attributes') {
            return this.csvEscape((row[knownIdx[key]] || '').trim());
          }
          const attribs = {};
          extraIdx.forEach((i) => {
            const attrKey = this.attribKeyFromHeader(this.pasteHeaders[i]);
            const val = String(row[i] || '').trim();
            if (attrKey && val) attribs[attrKey] = val;
          });
          return this.csvEscape(JSON.stringify(attribs));
        });
        lines.push(cells.join(','));
      });
      return `${lines.join('\n')}\n`;
    },

    guessMappedField(header) {
      const key = String(header || '').toLowerCase().replace(/[\s_-]/g, '');
      const aliases = {
        email: 'email',
        emailaddress: 'email',
        emailid: 'email',
        mail: 'email',
        name: 'name',
        fullname: 'name',
        firstname: 'firstname',
        first: 'firstname',
        fname: 'firstname',
        lastname: 'lastname',
        last: 'lastname',
        lname: 'lastname',
        surname: 'lastname',
        company: 'company',
        companyname: 'company',
        organization: 'company',
        organisation: 'company',
        attributes: 'attributes',
        attribs: 'attributes',
      };
      return aliases[key] || '';
    },

    parseImportFile(text) {
      this.fileRawText = text;
      const lines = text.split(/\r?\n/).filter((l) => l.trim().length > 0);
      if (lines.length === 0) {
        this.fileParseError = 'This file looks empty.';
        return;
      }
      const delim = this.detectDelim(lines[0]);
      this.form.delim = delim;
      const headers = this.parseCsvLine(lines[0], delim).filter((h) => h.length > 0);
      if (headers.length === 0) {
        this.fileParseError = 'Could not find a header row in this file.';
        return;
      }
      this.fileColumns = headers.map((h) => ({
        original: h,
        mapped: this.guessMappedField(h),
      }));
      this.fileRowCount = Math.max(0, lines.length - 1);
    },

    confirmMapping() {
      if (this.isZipFile) {
        this.fileStep = 3;
        return;
      }
      if (!this.fileColumns.some((c) => c.mapped === 'email')) {
        this.$utils.toast('Map at least one column to Email to continue.', 'is-danger');
        return;
      }
      const mapped = this.fileColumns.map((c) => c.mapped || c.original);
      const orig = this.fileColumns.map((c) => c.original);
      const changed = mapped.some((h, i) => h !== orig[i]);
      if (changed && this.fileRawText) {
        const nl = this.fileRawText.indexOf('\n');
        const rest = nl === -1 ? '' : this.fileRawText.slice(nl);
        const headerLine = mapped.map((h) => {
          if (/[",\n]/.test(h)) return `"${h.replace(/"/g, '""')}"`;
          return h;
        }).join(this.form.delim);
        const next = `${headerLine}${rest || '\n'}`;
        this.form.file = new File([next], this.form.file.name, { type: 'text/csv' });
        this.fileRawText = next;
      }
      this.fileStep = 3;
    },

    downloadExample() {
      const blob = new Blob([this.example], { type: 'text/csv;charset=utf-8' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'contacts-example.csv';
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    },

    clearFile() {
      this.resetFileWizard();
    },

    isFree() {
      if (this.status.status === 'none') {
        return true;
      }
      return false;
    },

    isRunning() {
      if (this.status.status === 'importing'
        || this.status.status === 'stopping') {
        return true;
      }
      return false;
    },

    isSuccessful() {
      return this.status.status === 'finished';
    },

    isFailed() {
      return (
        this.status.status === 'stopped'
        || this.status.status === 'failed'
      );
    },

    isDone() {
      if (this.status.status === 'finished'
        || this.status.status === 'stopped'
        || this.status.status === 'failed'
      ) {
        return true;
      }
      return false;
    },

    pollStatus() {
      clearInterval(this.pollID);

      this.pollID = setInterval(() => {
        this.$api.getImportStatus().then((data) => {
          this.isProcessing = false;
          this.isLoading = false;
          this.status = data;
          this.getLogs();

          if (!this.isRunning()) {
            clearInterval(this.pollID);
          }
        }, () => {
          this.isProcessing = false;
          this.isLoading = false;
          this.status = { status: 'none' };
          clearInterval(this.pollID);
        });
        return true;
      }, 250);
    },

    getLogs() {
      this.$api.getImportLogs().then((data) => {
        this.logs = data.split('\n').map((line) => line.replace(/\s+importer\.go:\d+:\s*/, ' *: '));
        Vue.nextTick(() => {
          const ref = document.getElementById('import-log');
          if (ref) {
            ref.scrollTop = ref.scrollHeight;
          }
        });
      });
    },

    stopImport() {
      this.isProcessing = true;
      this.$api.stopImport().then(() => {
        this.pollStatus();
        this.form.file = null;
        this.form.paste = '';
        this.clearPasteGrid();
        this.resetFileWizard();
        this.step = 'chooser';
      });
    },

    renderExample() {
      const h = 'email,firstname,lastname,company,attributes\n'
        + 'user1@mail.com,User,One,Acme Corp,"{""age"": 42, ""planet"": ""Mars""}"\n'
        + 'user2@mail.com,User,Two,Example Inc,"{""age"": 24, ""job"": ""Time Traveller""}"';

      this.example = h;
    },

    resetForm() {
      this.form.mode = 'subscribe';
      this.form.overwriteUserInfo = false;
      this.form.overwriteSubStatus = false;
      this.form.file = null;
      this.form.paste = '';
      this.form.lists = [];
      this.form.subStatus = 'unconfirmed';
      this.form.delim = ',';
      this.clearPasteGrid();
      this.resetFileWizard();
      this.preselectListFromQuery();
    },

    checkPasteData() {
      if (!this.pasteHeaders.length) {
        const parsed = this.parsePastedText(this.form.paste);
        if (!parsed) {
          this.$utils.toast('Please paste at least a header and one row.', 'is-danger');
          return;
        }
        this.applyPasteGrid(parsed);
      }
      const hasEmail = this.pasteHeaders.some((h) => this.guessMappedField(h) === 'email');
      if (!hasEmail) {
        this.$utils.toast('An Email column is required.', 'is-danger');
        return;
      }
      this.pasteChecked = true;
      this.pasteSummary = `${this.pasteRows.length} row(s) detected with ${this.pasteHeaders.length} column(s).`;
    },

    preselectListFromQuery() {
      const ids = this.$utils.parseQueryIDs(this.$route.query.list_id);
      if (ids.length === 0) return;

      const apply = () => {
        if (!this.lists.results) return;
        this.form.lists = this.lists.results.filter((l) => ids.indexOf(l.id) > -1);
      };

      if (this.lists.results && this.lists.results.length) {
        apply();
      } else {
        this.$api.getLists({ minimal: true, per_page: 'all', status: 'active' }).then(apply);
      }
    },

    onUpload() {
      if (this.step === 'paste') {
        if (!this.pasteChecked) return;
        const csv = this.buildImportCsvFromPaste();
        if (!csv) {
          this.$utils.toast('An Email column is required.', 'is-danger');
          return;
        }
        this.form.delim = ',';
        this.form.file = new File([csv], 'paste-import.csv', { type: 'text/csv' });
      }

      if (this.form.mode === 'subscribe' && this.form.overwriteSubStatus) {
        this.$utils.confirm(this.$t('import.subscribeWarning'), this.onSubmit, this.resetForm);
        return;
      }

      this.onSubmit();
    },

    onSubmit() {
      this.isProcessing = true;

      const params = new FormData();
      params.set('params', JSON.stringify({
        mode: this.form.mode,
        subscription_status: this.form.subStatus,
        delim: this.form.delim,
        lists: this.form.lists.map((l) => l.id),
        overwrite_userinfo: this.form.overwriteUserInfo,
        overwrite_subscription_status: this.form.overwriteSubStatus,
      }));
      params.set('file', this.form.file);

      this.$api.importSubscribers(params).then(() => {
        this.$utils.toast(this.$t('import.importStarted'));
        this.pollStatus();
      }, () => {
        this.isProcessing = false;
        this.form.file = null;
      });
    },
  },

  computed: {
    ...mapState(['lists']),

    canSubmitFile() {
      if (this.fileStep < 3) return false;
      if (this.form.mode === 'subscribe' && this.form.lists.length === 0) return false;
      return !!this.form.file;
    },

    canSubmitPaste() {
      if (!this.pasteChecked) return false;
      if (this.form.mode === 'subscribe' && this.form.lists.length === 0) return false;
      return this.pasteRows.length > 0 && this.pasteHeaders.some((h) => this.guessMappedField(h) === 'email');
    },

    hasPasteData() {
      return this.pasteHeaders.length > 0 || !!(this.form.paste || '').trim();
    },

    pastePreviewRows() {
      return (this.pasteRows || []).slice(0, 25);
    },

    progress() {
      if (!this.status || !this.status.total > 0) {
        return 0;
      }
      return Math.ceil((this.status.imported / this.status.total) * 100);
    },
  },

  mounted() {
    this.renderExample();
    this.pollStatus();
    this.preselectListFromQuery();
  },
});
</script>
