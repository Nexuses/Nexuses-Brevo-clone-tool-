<template>
  <section class="crm-page templates-page">
    <marketing-subnav />

    <header class="columns page-header">
      <div class="column is-8">
        <h1 class="title is-4">Templates</h1>
      </div>
      <div class="column has-text-right">
        <b-button
          v-if="$can('templates:manage')"
          type="is-dark"
          icon-left="plus"
          class="btn-new"
          data-cy="btn-new"
          @click="createTemplate"
        >
          Create Template
        </b-button>
      </div>
    </header>

    <div class="tpl-tabs">
      <button type="button" class="tpl-tabs__tab is-active">Email</button>
      <button type="button" class="tpl-tabs__tab" disabled>WhatsApp</button>
    </div>

    <div class="tpl-toolbar">
      <b-field>
        <b-input
          v-model="query"
          expanded
          placeholder="Search for templates"
          icon="magnify"
          data-cy="template-search"
        />
      </b-field>
      <span class="tpl-toolbar__meta">
        {{ filteredTemplates.length }} of {{ templateList.length }}
      </span>
    </div>

    <b-loading :active="loading.templates" :is-full-page="false" />

    <div class="tpl-list">
      <div
        v-for="t in filteredTemplates"
        :key="t.id"
        class="tpl-card"
      >
        <div class="tpl-card__main">
          <a href="#" class="tpl-card__title" @click.prevent="editTemplate(t)">
            {{ t.name }}
          </a>
          <div class="tpl-card__meta">
            #{{ t.id }} · Last edited on {{ formatDate(t.updatedAt || t.createdAt) }}
          </div>
          <div class="tpl-card__status" :class="statusClass(t)">
            <span class="dot" />
            {{ statusLabel(t) }}
          </div>
        </div>

        <div class="tpl-card__actions">
          <a
            href="#"
            class="tpl-card__edit"
            aria-label="Edit"
            @click.prevent="editTemplate(t)"
          >
            <b-tooltip label="Edit" type="is-dark" position="is-top">
              <b-icon icon="pencil-outline" />
            </b-tooltip>
          </a>

          <b-dropdown position="is-bottom-left" class="campaign-actions-menu">
            <template #trigger>
              <button type="button" class="campaign-actions-trigger" aria-label="Actions">
                <span class="campaign-kebab" aria-hidden="true"><span /><span /><span /></span>
              </button>
            </template>
            <div class="campaign-actions-panel">
              <a
                href="#"
                class="campaign-action"
                aria-label="Preview"
                @click.prevent="previewTemplate(t)"
              >
                <b-tooltip :label="$t('templates.preview')" type="is-dark" position="is-left">
                  <b-icon icon="file-find-outline" />
                </b-tooltip>
              </a>
              <a
                href="#"
                class="campaign-action"
                aria-label="Clone"
                @click.prevent="$utils.prompt($t('globals.buttons.clone'),
                  { placeholder: $t('globals.fields.name'), value: $t('campaigns.copyOf', { name: t.name }) },
                  (name) => cloneTemplate(name, t))"
              >
                <b-tooltip :label="$t('globals.buttons.clone')" type="is-dark" position="is-left">
                  <b-icon icon="file-multiple-outline" />
                </b-tooltip>
              </a>
              <a
                v-if="!t.isDefault && t.type === 'campaign'"
                href="#"
                class="campaign-action"
                aria-label="Make default"
                @click.prevent="$utils.confirm(null, () => makeTemplateDefault(t))"
              >
                <b-tooltip :label="$t('templates.makeDefault')" type="is-dark" position="is-left">
                  <b-icon icon="check-circle-outline" />
                </b-tooltip>
              </a>
              <a
                v-if="!t.isDefault"
                href="#"
                class="campaign-action"
                aria-label="Delete"
                @click.prevent="$utils.confirm(null, () => deleteTemplate(t))"
              >
                <b-tooltip :label="$t('globals.buttons.delete')" type="is-dark" position="is-left">
                  <b-icon icon="trash-can-outline" />
                </b-tooltip>
              </a>
            </div>
          </b-dropdown>
        </div>
      </div>

      <div v-if="!loading.templates && filteredTemplates.length === 0" class="tpl-empty">
        No templates found.
      </div>
    </div>

    <campaign-preview
      v-if="previewItem"
      type="template"
      :id="previewItem.id"
      :template-type="previewItem.type"
      :title="previewItem.name"
      @close="closePreview"
    />
  </section>
</template>

<script>
import dayjs from 'dayjs';
import Vue from 'vue';
import { mapState } from 'vuex';
import CampaignPreview from '../components/CampaignPreview.vue';
import MarketingSubnav from '../components/MarketingSubnav.vue';

export default Vue.extend({
  components: {
    CampaignPreview,
    MarketingSubnav,
  },

  data() {
    return {
      query: '',
      previewItem: null,
    };
  },

  computed: {
    ...mapState(['templates', 'loading']),

    templateList() {
      return Array.isArray(this.templates) ? this.templates : [];
    },

    filteredTemplates() {
      const q = this.query.trim().toLowerCase();
      if (!q) return this.templateList;
      return this.templateList.filter((t) => (
        (t.name || '').toLowerCase().includes(q)
        || String(t.id).includes(q)
      ));
    },
  },

  methods: {
    fetchTemplates() {
      this.$api.getTemplates();
    },

    formatDate(iso) {
      return dayjs(iso).format('MMM D, YYYY h:mm A');
    },

    statusLabel(t) {
      if (t.isDefault) return 'Active · Default';
      return 'Active';
    },

    statusClass(t) {
      return t.isDefault ? 'is-active' : 'is-active';
    },

    createTemplate() {
      this.$router.push({ name: 'template', params: { id: 'new' } });
    },

    editTemplate(t) {
      this.$router.push({ name: 'template', params: { id: t.id } });
    },

    previewTemplate(t) {
      this.previewItem = t;
    },

    closePreview() {
      this.previewItem = null;
    },

    cloneTemplate(name, t) {
      const data = {
        name,
        type: t.type,
        subject: t.subject,
        body: t.body,
        body_source: t.bodySource,
      };
      this.$api.createTemplate(data).then((d) => {
        this.$api.getTemplates();
        this.$utils.toast(`'${d.name}' created`);
      });
    },

    makeTemplateDefault(tpl) {
      this.$api.makeTemplateDefault(tpl.id).then(() => {
        this.$api.getTemplates();
        this.$utils.toast(this.$t('globals.messages.created', { name: tpl.name }));
      });
    },

    deleteTemplate(tpl) {
      this.$api.deleteTemplate(tpl.id).then(() => {
        this.$api.getTemplates();
        this.$utils.toast(this.$t('globals.messages.deleted', { name: tpl.name }));
      });
    },
  },

  created() {
    this.$root.$on('page.refresh', this.fetchTemplates);
  },

  destroyed() {
    this.$root.$off('page.refresh', this.fetchTemplates);
  },

  mounted() {
    this.$api.getTemplates();
  },
});
</script>
