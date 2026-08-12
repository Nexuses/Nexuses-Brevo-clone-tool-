<template>
  <section class="templates templates-brevo bv-page">
    <marketing-subnav />

    <header class="templates-brevo__header">
      <div class="templates-brevo__header-main">
        <h1 class="templates-brevo__title">
          Templates
          <span v-if="templates.length > 0" class="has-text-grey-light">({{ templates.length }})</span>
        </h1>
        <p class="templates-brevo__lead">
          Create and manage reusable email templates for your campaigns.
        </p>
        <div class="templates-brevo__links">
          <a
            href="https://listmonk.app/docs/templating/"
            target="_blank"
            rel="noopener noreferrer"
            class="templates-brevo__link"
          >
            Get started with templates
            <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
              <path d="M3.5 3.5h5v5M8.5 3.5L3.5 8.5" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
          </a>
        </div>
      </div>
      <button
        v-if="$can('templates:manage')"
        type="button"
        class="templates-brevo__create"
        data-cy="btn-new"
        @click="showNewForm"
      >
        <span class="templates-brevo__create-plus" aria-hidden="true">+</span>
        Create template
      </button>
    </header>

    <b-table
      :data="templates"
      :hoverable="true"
      :loading="loading.templates"
      default-sort="createdAt"
      class="templates-brevo__table"
    >
      <b-table-column v-slot="props" field="name" label="Template" :td-attrs="$utils.tdID" sortable>
        <a href="#" class="templates-brevo__name" @click.prevent="showEditForm(props.row)">
          {{ props.row.name }}
        </a>
        <b-tag v-if="props.row.isDefault" class="ml-2">
          {{ $t('templates.default') }}
        </b-tag>

        <p class="is-size-7 has-text-grey" v-if="props.row.type === 'tx'">
          {{ props.row.subject }}
        </p>
      </b-table-column>

      <b-table-column v-slot="props" field="type" :label="$t('globals.fields.type')" sortable>
        <b-tag v-if="props.row.type === 'campaign'" :class="props.row.type" :data-cy="`type-${props.row.type}`">
          {{ $tc('templates.typeCampaignHTML') }}
        </b-tag>
        <b-tag v-else-if="props.row.type === 'campaign_visual'" :class="props.row.type"
          :data-cy="`type-${props.row.type}`">
          {{ $tc('templates.typeCampaignVisual') }}
        </b-tag>
        <b-tag v-else-if="props.row.type === 'campaign_maily'" :class="props.row.type"
          :data-cy="`type-${props.row.type}`">
          {{ $tc('templates.typeCampaignMaily') }}
        </b-tag>
        <b-tag v-else :class="props.row.type" :data-cy="`type-${props.row.type}`">
          {{ $tc('templates.typeTransactional') }}
        </b-tag>
      </b-table-column>

      <b-table-column v-slot="props" field="id" :label="$t('globals.fields.id')" sortable>
        #{{ props.row.id }}
      </b-table-column>

      <b-table-column v-slot="props" field="createdAt" :label="$t('globals.fields.createdAt')" sortable>
        {{ $utils.niceDate(props.row.createdAt) }}
      </b-table-column>

      <b-table-column v-slot="props" field="updatedAt" :label="$t('globals.fields.updatedAt')" sortable>
        {{ $utils.niceDate(props.row.updatedAt) }}
      </b-table-column>

      <b-table-column v-slot="props" cell-class="actions" align="right" label="Actions">
        <b-dropdown position="is-bottom-left">
          <template #trigger>
            <button type="button" class="templates-brevo__kebab" aria-label="Actions">
              <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
                <circle cx="8" cy="3.5" r="1.4" fill="currentColor" />
                <circle cx="8" cy="8" r="1.4" fill="currentColor" />
                <circle cx="8" cy="12.5" r="1.4" fill="currentColor" />
              </svg>
            </button>
          </template>
          <b-dropdown-item @click="previewTemplate(props.row)">
            {{ $t('templates.preview') }}
          </b-dropdown-item>
          <b-dropdown-item @click="showEditForm(props.row)">
            {{ $t('globals.buttons.edit') }}
          </b-dropdown-item>
          <b-dropdown-item
            @click="$utils.prompt(`Clone template`,
              { placeholder: 'Name', value: `Copy of ${props.row.name}` },
              (name) => cloneTemplate(name, props.row))"
          >
            {{ $t('globals.buttons.clone') }}
          </b-dropdown-item>
          <b-dropdown-item
            v-if="!props.row.isDefault && props.row.type === 'campaign'"
            @click="$utils.confirm(null, () => makeTemplateDefault(props.row))"
          >
            {{ $t('templates.makeDefault') }}
          </b-dropdown-item>
          <b-dropdown-item
            v-if="!props.row.isDefault"
            class="has-text-danger"
            @click="$utils.confirm(null, () => deleteTemplate(props.row))"
          >
            {{ $t('globals.buttons.delete') }}
          </b-dropdown-item>
        </b-dropdown>
      </b-table-column>

      <template #empty v-if="!loading.templates">
        <empty-placeholder />
      </template>
    </b-table>

    <!-- Add / edit form modal -->
    <b-modal scroll="keep" :aria-modal="true" :active.sync="isFormVisible" :width="1200" :can-cancel="false"
      class="template-modal">
      <template-form :data="curItem" :is-editing="isEditing" @finished="formFinished" />
    </b-modal>

    <campaign-preview v-if="previewItem" type="template" :id="previewItem.id" :template-type="previewItem.type"
      :title="previewItem.name" @close="closePreview" />
  </section>
</template>

<script>
import Vue from 'vue';
import { mapState } from 'vuex';
import CampaignPreview from '../components/CampaignPreview.vue';
import EmptyPlaceholder from '../components/EmptyPlaceholder.vue';
import MarketingSubnav from '../components/MarketingSubnav.vue';

import TemplateForm from './TemplateForm.vue';

export default Vue.extend({
  components: {
    CampaignPreview,
    TemplateForm,
    EmptyPlaceholder,
    MarketingSubnav,
  },

  data() {
    return {
      curItem: null,
      isEditing: false,
      isFormVisible: false,
      previewItem: null,
    };
  },

  methods: {
    fetchTemplates() {
      this.$api.getTemplates();
    },

    // Show the edit form.
    showEditForm(data) {
      this.curItem = data;
      this.isFormVisible = true;
      this.isEditing = true;
    },

    // Show the new form.
    showNewForm() {
      this.curItem = { type: 'campaign' };
      this.isFormVisible = true;
      this.isEditing = false;
    },

    formFinished() {
      this.$api.getTemplates();
    },

    previewTemplate(c) {
      this.previewItem = c;
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
        this.$emit('finished');
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

  computed: {
    ...mapState(['templates', 'loading']),
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
