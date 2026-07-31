<template>
  <div class="maily-editor-wrapper">
    <iframe ref="mailyEditor" id="maily-editor" class="maily-editor" title="Maily email editor" />
  </div>
</template>

<script>
export default {
  props: {
    source: { type: String, default: '' },
    height: { type: String, default: 'auto' },
  },

  methods: {
    loadAssets() {
      return new Promise((resolve, reject) => {
        const iframe = this.$refs.mailyEditor;
        const doc = iframe.contentDocument;
        if (iframe.contentWindow.MailyBuilder) {
          resolve();
          return;
        }

        const link = doc.createElement('link');
        link.rel = 'stylesheet';
        link.href = '/admin/static/maily-builder/maily-builder.css';
        doc.head.appendChild(link);

        const script = doc.createElement('script');
        script.id = 'maily-builder-script';
        script.src = '/admin/static/maily-builder/maily-builder.umd.js';
        script.onload = () => resolve();
        script.onerror = reject;
        doc.head.appendChild(script);
      });
    },

    render(source) {
      const iframe = this.$refs.mailyEditor;
      const mb = iframe.contentWindow.MailyBuilder;
      if (!mb) {
        return;
      }

      let data = null;
      if (source) {
        try {
          data = typeof source === 'string' ? JSON.parse(source) : source;
        } catch (e) {
          /* eslint-disable-next-line no-console */
          console.error('Invalid Maily JSON:', e);
        }
      }

      const onChange = (json, body) => {
        const tpl = body.replace(/\{\{[^}]*\}\}/g, (match) => match.replace(/&quot;/g, '"'));
        this.$emit('change', { source: JSON.stringify(json), body: tpl });
      };

      if (!mb.isRendered('maily-editor-container')) {
        mb.render('maily-editor-container', { data, onChange });
      } else if (data && mb.setContent) {
        mb.setContent('maily-editor-container', data, onChange);
      }
    },
  },

  mounted() {
    const iframe = this.$refs.mailyEditor;
    iframe.style.height = this.height;

    iframe.srcdoc = `
      <!DOCTYPE html>
      <html>
        <head>
          <style>
            html, body { margin: 0; padding: 0; height: 100%; }
            #maily-editor-container { width: 100%; min-height: 100%; }
          </style>
        </head>
        <body>
          <div id="maily-editor-container"></div>
        </body>
      </html>
    `;

    iframe.onload = () => {
      this.loadAssets().then(() => {
        let source = null;
        if (this.$props.source) {
          source = this.$props.source;
        }
        this.render(source);
      }).catch((error) => {
        /* eslint-disable-next-line no-console */
        console.error('Failed to load Maily editor:', error);
      });
    };
  },
};
</script>

<style lang="css">
.maily-editor-wrapper {
  width: 100%;
  border: 1px solid #eaeaea;
  max-width: 100vw;
}

#maily-editor {
  position: relative;
  border: none;
  width: 100%;
  min-height: 500px;
}
</style>
