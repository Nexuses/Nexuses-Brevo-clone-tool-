<template>
  <section class="template-editor">
    <marketing-subnav />

    <header class="columns page-header">
      <div class="column is-6">
        <p class="mb-2">
          <router-link class="crm-link" :to="{ name: 'templates' }">
            ← Back to templates
          </router-link>
        </p>
        <h1 class="title is-4">
          {{ isEditing ? (template.name || 'Edit template') : $t('templates.newTemplate') }}
        </h1>
        <p v-if="isEditing && template.id" class="has-text-grey is-size-7">
          {{ $t('globals.fields.id') }}: #{{ template.id }}
        </p>
      </div>
    </header>

    <b-loading :active="pageLoading" :is-full-page="false" />

    <template-form
      v-if="ready"
      :key="formKey"
      :data="template"
      :is-editing="isEditing"
      as-page
      @finished="onFinished"
      @cancel="goBack"
    />
  </section>
</template>

<script>
import MarketingSubnav from '../components/MarketingSubnav.vue';
import TemplateForm from './TemplateForm.vue';

export default {
  name: 'TemplateEditor',
  components: {
    MarketingSubnav,
    TemplateForm,
  },

  data() {
    return {
      pageLoading: false,
      ready: false,
      template: { type: 'campaign', name: '', body: '' },
    };
  },

  computed: {
    isEditing() {
      return this.$route.params.id !== 'new';
    },
    formKey() {
      return `${this.$route.params.id}-${this.template.id || 'new'}`;
    },
  },

  methods: {
    goBack() {
      this.$router.push({ name: 'templates' });
    },

    onFinished(tpl) {
      if (tpl && tpl.id && !this.isEditing) {
        this.$router.replace({ name: 'template', params: { id: tpl.id } });
        return;
      }
      this.$api.getTemplates();
    },

    async load() {
      this.ready = false;
      if (!this.isEditing) {
        this.template = { type: 'campaign', name: '', body: '' };
        this.ready = true;
        return;
      }

      this.pageLoading = true;
      try {
        const data = await this.$api.getTemplate(this.$route.params.id);
        this.template = data || { type: 'campaign', name: '', body: '' };
        this.ready = true;
      } catch (e) {
        this.$utils.toast('Template not found');
        this.goBack();
      } finally {
        this.pageLoading = false;
      }
    },
  },

  watch: {
    // Reload when navigating between new / edit template routes.
    '$route.params.id': function onRouteIdChange() {
      this.load();
    },
  },

  mounted() {
    this.load();
  },
};
</script>
