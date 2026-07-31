<template>
  <!-- Two-way Data-Binding -->
  <section class="editor">
    <div class="columns">
      <div class="column is-three-quarters is-inline-flex">
        <b-field :label="$t('campaigns.format')" label-position="on-border" class="mr-4 mb-0">
          <b-select v-model="contentTypeSel" :disabled="disabled" name="content_type">
            <option v-for="(name, f) in contentTypes" :key="f" name="format" :value="f" :data-cy="`check-${f}`">
              {{ name }}
            </option>
          </b-select>
        </b-field>

        <b-field v-if="!isBlockEditor" :label="$tc('globals.terms.template')" label-position="on-border">
          <b-select :placeholder="$t('globals.terms.none')" v-model="templateId" name="template" :disabled="disabled">
            <template v-for="t in validTemplates">
              <option :value="t.id" :key="t.id">
                {{ t.name }}
              </option>
            </template>
          </b-select>
        </b-field>

        <div v-else>
          <b-button v-if="!isBlockTplSelector" @click="onShowBlockTplSelector" type="is-ghost"
            icon-left="file-find-outline" data-cy="btn-select-block-tpl">
            {{ blockImportLabel }}
          </b-button>
          <b-field v-else :label="$tc('globals.terms.template')" label-position="on-border">
            <b-select :placeholder="$t('globals.terms.none')" v-model="blockTemplateId"
              @input="() => isBlockTplDisabled = false" name="template" :disabled="disabled"
              class="copy-visual-template-list">
              <template v-for="t in validTemplates">
                <option :value="t.id" :key="t.id">
                  {{ t.name }}
                </option>
              </template>
            </b-select>

            <b-button :disabled="disabled || isBlockTplDisabled || !blockTemplateId" class="ml-3"
              @click="onImportBlockTpl" type="is-primary" icon-left="content-save-outline"
              data-cy="btn-save-block-tpl">
              {{ $t('globals.terms.import') }}

              <span class="spinner is-tiny" v-if="loading.templates">
                <b-loading :is-full-page="false" active />
              </span>
            </b-button>
          </b-field>
        </div>
      </div>
      <div class="column is- has-text-right">
        <b-button @click="onTogglePreview" type="is-primary" icon-left="file-find-outline" data-cy="btn-preview"
          aria-keyshortcuts="F9">
          <span class="has-kbd">{{ $t('campaigns.preview') }} <span class="kbd">F9</span></span>
        </b-button>
      </div>
    </div>

    <!-- wsywig //-->
    <richtext-editor v-if="self.contentType === 'richtext'" :disabled="disabled" v-model="self.body" />

    <!-- visual editor //-->
    <visual-editor v-if="self.contentType === 'visual'" :source="self.bodySource" @change="onBlockEditorChange"
      height="65vh" ref="visualEditor" />

    <!-- maily editor //-->
    <maily-editor v-if="self.contentType === 'maily'" :source="self.bodySource" @change="onBlockEditorChange"
      height="65vh" ref="mailyEditor" />

    <!-- raw html editor //-->
    <code-editor lang="html" v-if="self.contentType === 'html'" v-model="self.body" key="editor-html" />

    <!-- markdown editor //-->
    <code-editor lang="markdown" v-if="self.contentType === 'markdown'" v-model="self.body" key="editor-markdown" />

    <!-- plain text //-->
    <b-input v-if="self.contentType === 'plain'" v-model="self.body" type="textarea" name="content" ref="plainEditor"
      class="plain-editor" />

    <!-- campaign preview //-->
    <campaign-preview v-if="isPreviewing" is-post @close="onTogglePreview" type="campaign" :id="id" :title="title"
      :content-type="self.contentType" :template-id="templateId" :body="self.body" />
  </section>
</template>

<script>
import { html as beautifyHTML } from 'js-beautify';
import TurndownService from 'turndown';
import { mapState } from 'vuex';

import CampaignPreview from './CampaignPreview.vue';
import VisualEditor from './VisualEditor.vue';
import MailyEditor from './MailyEditor.vue';
import RichtextEditor from './RichtextEditor.vue';
import markdownToVisualBlock from './editor';
import CodeEditor from './CodeEditor.vue';

const turndown = new TurndownService();

const DEFAULT_MAILY_DOC = {
  type: 'doc',
  content: [
    {
      type: 'paragraph',
      attrs: { textAlign: 'left' },
      content: [{ type: 'text', text: '' }],
    },
  ],
};

export default {
  components: {
    CampaignPreview,
    'code-editor': CodeEditor,
    'visual-editor': VisualEditor,
    'maily-editor': MailyEditor,
    'richtext-editor': RichtextEditor,
  },

  props: {
    contentTypes: { type: Object, default: () => ({}) },
    id: { type: Number, default: 0 },
    title: { type: String, default: '' },
    disabled: { type: Boolean, default: false },
    templates: { type: Array, default: null },

    value: {
      type: Object,
      default: () => ({
        body: '',
        bodySource: null,
        contentType: '',
        templateId: null,
      }),
    },
  },

  data() {
    return {
      isPreviewing: false,
      isBlockTplSelector: false,
      isBlockTplDisabled: false,
      contentTypeSel: this.$props.value.contentType,
      templateId: null,
      blockTemplateId: null,
    };
  },

  methods: {
    onContentTypeChange(to, from) {
      if (!this.self.body.trim()) {
        this.convertContentType(to, from);
        return;
      }

      this.$utils.confirm(
        this.$t('campaigns.confirmSwitchFormat'),
        () => {
          this.convertContentType(to, from);
        },
        () => {
          this.contentTypeSel = from;
        },
      );
    },

    convertContentType(to, from) {
      let body = this.self.body ?? '';
      let bodySource = null;

      let skip = false;

      let isHTML = false;
      if (from === 'richtext' || from === 'html' || from === 'visual' || from === 'maily') {
        const d = document.createElement('div');
        d.innerHTML = body;
        body = this.beautifyHTML(d.innerHTML.trim());
        isHTML = true;
      }

      if (isHTML) {
        switch (to) {
          case 'plain': {
            const d = document.createElement('div');
            d.innerHTML = body;
            body = this.trimLines(d.innerText.trim(), true);
            break;
          }

          case 'markdown': {
            body = turndown.turndown(body).replace(/\n\n+/ig, '\n\n');
            break;
          }

          case 'visual': {
            const md = turndown.turndown(body).replace(/\n\n+/ig, '\n\n');
            bodySource = JSON.stringify(markdownToVisualBlock(md));
            break;
          }

          case 'maily': {
            bodySource = JSON.stringify(DEFAULT_MAILY_DOC);
            body = '';
            break;
          }

          default:
            break;
        }
      } else if (from === 'markdown' && (to === 'richtext' || to === 'html')) {
        skip = true;
        this.$api.convertCampaignContent({
          id: 1, body, from, to,
        }).then((data) => {
          this.$nextTick(() => {
            this.self.contentType = to;
            this.self.body = this.beautifyHTML(data.trim());
          });
        });
      } else if (from === 'plain' && (to === 'richtext' || to === 'html')) {
        body = body.replace(/\n/ig, '<br>\n');
      } else if (to === 'visual') {
        bodySource = JSON.stringify(markdownToVisualBlock(body));
      } else if (to === 'maily') {
        bodySource = JSON.stringify(DEFAULT_MAILY_DOC);
        body = '';
      }

      if (to === 'visual' || from === 'visual' || to === 'maily' || from === 'maily') {
        this.templateId = null;
        this.self.templateId = null;
      }

      if (!skip) {
        this.$nextTick(() => {
          this.self.contentType = to;
          this.self.body = body;
          this.self.bodySource = bodySource;
        });
      }
    },

    onTogglePreview() {
      this.isPreviewing = !this.isPreviewing;
    },

    onKeyboardShortcut(e) {
      if (e.key === 'F9') {
        this.onTogglePreview();
        e.preventDefault();
      }

      if (e.ctrlKey && e.key === 's') {
        this.$events.$emit('campaign.update');
        e.preventDefault();
      }
    },

    onBlockEditorChange({ body, source }) {
      this.self.body = body;
      this.self.bodySource = source;
    },

    beautifyHTML(str) {
      let s = this.trimLines(str.replace(/(<(?!(\/)?a|span)([^>]+)>)/ig, '\n$1\n'), true);
      s = s.replace(/\n+/g, '\n');

      return beautifyHTML(s, {
        indent_size: 4,
        indent_char: ' ',
        max_preserve_newlines: 2,
        inline: ['h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'b', 'strong', 'span', 'em', 'i', 'code', 'a'],
      }).trim();
    },

    trimLines(str, removeEmptyLines) {
      const out = str.split('\n');
      for (let i = 0; i < out.length; i += 1) {
        const line = out[i].trim();
        if (removeEmptyLines) {
          out[i] = line;
        } else if (line === '') {
          out[i] = '';
        }
      }

      return out.join('\n').replace(/\n\s*\n\s*\n/g, '\n\n');
    },

    onShowBlockTplSelector() {
      this.isBlockTplSelector = true;
      this.setDefaultTemplate();
    },

    onImportBlockTpl() {
      if (!this.blockTemplateId) {
        return;
      }

      this.$utils.confirm(
        this.$t('campaigns.confirmOverwriteContent'),
        () => {
          this.$api.getTemplate(this.blockTemplateId).then((data) => {
            this.self.body = data.body;
            this.self.bodySource = data.bodySource;
            this.isBlockTplDisabled = true;

            const parsed = JSON.parse(data.bodySource);
            if (this.self.contentType === 'maily' && this.$refs.mailyEditor) {
              this.$refs.mailyEditor.render(parsed);
            } else if (this.$refs.visualEditor) {
              this.$refs.visualEditor.render(parsed);
            }
          });
        },
      );
    },

    setDefaultTemplate() {
      if (this.isBlockEditor) {
        this.blockTemplateId = this.validTemplates[0]?.id || null;
      } else {
        if (this.templateId) {
          return;
        }

        const defaultTemplate = this.validTemplates.find((t) => t.isDefault === true);
        this.templateId = defaultTemplate?.id || this.validTemplates[0]?.id || null;
      }
    },
  },

  mounted() {
    this.contentTypeSel = this.value.contentType;
    this.templateId = this.value.templateId;

    window.addEventListener('keydown', this.onKeyboardShortcut);

    this.$events.$on('campaign.preview', () => {
      this.isPreviewing = true;
    });
  },

  beforeDestroy() {
    window.removeEventListener('keydown', this.onKeyboardShortcut);
    this.$events.$off('campaign.preview');
  },

  computed: {
    ...mapState(['serverConfig', 'loading']),

    self: {
      get() {
        return this.value;
      },
      set(val) {
        this.$emit('input', val);
      },
    },

    isBlockEditor() {
      return this.self.contentType === 'visual' || this.self.contentType === 'maily';
    },

    blockImportLabel() {
      if (this.self.contentType === 'maily') {
        return this.$t('campaigns.importMailyTemplate');
      }
      return this.$t('campaigns.importVisualTemplate');
    },

    validTemplates() {
      let typ = 'campaign';
      if (this.self.contentType === 'visual') {
        typ = 'campaign_visual';
      } else if (this.self.contentType === 'maily') {
        typ = 'campaign_maily';
      }
      return this.templates.filter((t) => (t.type === typ));
    },
  },

  watch: {
    validTemplates() {
      this.setDefaultTemplate();
    },

    contentTypeSel(to, from) {
      if (from !== to && to !== this.self.contentType) {
        this.onContentTypeChange(to, from);
      }
    },

    templateId(to) {
      if (this.self.templateId === to) {
        return;
      }

      this.self.templateId = to;
    },
  },
};

</script>
