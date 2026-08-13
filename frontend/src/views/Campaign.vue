<template>
  <section
    class="campaign bv-page"
    :class="{ 'is-design-open': isDesignOpen, 'is-report': showReport, 'is-setup': showSetup }"
  >
    <marketing-subnav v-if="!isDesignOpen && !showReport && !showSetup" />

    <campaign-report
      v-if="showReport"
      :campaign="data"
      @edit="goToEditor"
    />

    <campaign-setup
      v-else-if="showSetup"
      :form="form"
      :data="data"
      :is-new="isNew"
      :can-edit="canEdit"
      :can-manage="canManage"
      :can-send="canSend"
      :can-un-schedule="canUnSchedule"
      :loading="!!loading.campaigns"
      :lists="allLists"
      :templates="templates"
      :content-types="contentTypes"
      :from-email-default="serverConfig.from_email || ''"
      @save="onSetupSave"
      @ensure="onEnsureSaved"
      @schedule="onSetupSchedule"
      @unschedule="onUnschedule"
      @test="onSubmit('test')"
      @attach="isAttachModalOpen = true"
      @design-open="isDesignOpen = $event"
    />

    <b-modal scroll="keep" :aria-modal="true" :active.sync="isAttachModalOpen" :width="900">
      <div class="modal-card content" style="width: auto">
        <section expanded class="modal-card-body">
          <media is-modal @selected="onAttachSelect" />
        </section>
      </div>
    </b-modal>

    <campaign-preview v-if="isPreviewingArchive" @close="onToggleArchivePreview" type="campaign" :id="data.id"
      :archive-meta="form.archiveMetaStr" :title="data.title" :content-type="data.contentType"
      :template-id="form.archiveTemplateId" is-post is-archive />
  </section>
</template>

<script>
import dayjs from 'dayjs';
import htmlToPlainText from 'textversionjs';
import Vue from 'vue';
import { mapState } from 'vuex';

import CampaignPreview from '../components/CampaignPreview.vue';
import CampaignReport from '../components/CampaignReport.vue';
import CampaignSetup from '../components/CampaignSetup.vue';
import MarketingSubnav from '../components/MarketingSubnav.vue';
import Media from './Media.vue';

export default Vue.extend({
  components: {
    Media,
    CampaignPreview,
    CampaignReport,
    CampaignSetup,
    MarketingSubnav,
  },

  data() {
    return {
      contentTypes: Object.freeze({
        richtext: this.$t('campaigns.richText'),
        html: this.$t('campaigns.rawHTML'),
        markdown: this.$t('campaigns.markdown'),
        plain: this.$t('campaigns.plainText'),
        visual: this.$t('campaigns.visual'),
        maily: this.$t('campaigns.maily'),
      }),

      isNew: false,
      isEditing: false,
      isHeadersVisible: false,
      isAttachFieldVisible: false,
      isAttachModalOpen: false,
      isPreviewingArchive: false,
      isDesignOpen: false,
      activeTab: 'campaign',

      data: {},

      // IDs from ?list_id query param.
      selListIDs: [],

      // Binds form input values.
      form: {
        archiveSlug: null,
        name: '',
        subject: '',
        fromEmail: '',
        replyTo: '',
        headersStr: '[]',
        headers: [],
        attribsStr: '{}',
        messenger: 'email',
        lists: [],
        tags: [],
        sendAt: null,
        content: {
          contentType: 'richtext',
          body: '',
          bodySource: null,
          templateId: null,
        },
        altbody: null,
        media: [],

        // Parsed Date() version of send_at from the API.
        sendAtDate: null,
        sendLater: false,
        archive: false,
        archiveMetaStr: '{}',
        archiveMeta: {},
        testEmails: [],
        previewText: '',
        tracking: {
          enabled: true,
          track_views: true,
          track_clicks: true,
          notes: '',
          utm_source: '',
          utm_medium: '',
          utm_campaign: '',
          utm_content: '',
          utm_term: '',
          custom_params: {},
        },
        trackingIndividualMode: 'default',
        trackingCustomParamsStr: '{}',
      },
    };
  },

  methods: {
    parseTrackingConfig(raw) {
      const defaults = {
        enabled: true,
        track_views: true,
        track_clicks: true,
        notes: '',
        utm_source: '',
        utm_medium: '',
        utm_campaign: '',
        utm_content: '',
        utm_term: '',
        custom_params: {},
      };
      let tc = {};
      try {
        tc = raw && typeof raw === 'object' ? raw : JSON.parse(raw || '{}');
      } catch {
        tc = {};
      }
      const tracking = {
        ...defaults,
        enabled: tc.enabled ?? defaults.enabled,
        track_views: tc.track_views ?? tc.trackViews ?? defaults.track_views,
        track_clicks: tc.track_clicks ?? tc.trackClicks ?? defaults.track_clicks,
        notes: tc.notes ?? '',
        utm_source: tc.utm_source ?? tc.utmSource ?? '',
        utm_medium: tc.utm_medium ?? tc.utmMedium ?? '',
        utm_campaign: tc.utm_campaign ?? tc.utmCampaign ?? '',
        utm_content: tc.utm_content ?? tc.utmContent ?? '',
        utm_term: tc.utm_term ?? tc.utmTerm ?? '',
        custom_params: tc.custom_params ?? tc.customParams ?? {},
      };
      let trackingIndividualMode = 'default';
      if (tc.individual_tracking === true) {
        trackingIndividualMode = 'on';
      } else if (tc.individual_tracking === false) {
        trackingIndividualMode = 'off';
      }
      return {
        tracking,
        trackingIndividualMode,
        trackingCustomParamsStr: JSON.stringify(tracking.custom_params || {}, null, 2),
      };
    },

    buildTrackingConfig() {
      const t = this.form.tracking;
      const tc = {
        enabled: t.enabled,
        track_views: t.track_views,
        track_clicks: t.track_clicks,
        notes: t.notes || '',
        utm_source: t.utm_source || '',
        utm_medium: t.utm_medium || '',
        utm_campaign: t.utm_campaign || '',
        utm_content: t.utm_content || '',
        utm_term: t.utm_term || '',
      };
      if (this.form.trackingIndividualMode === 'on') {
        tc.individual_tracking = true;
      } else if (this.form.trackingIndividualMode === 'off') {
        tc.individual_tracking = false;
      } else {
        tc.individual_tracking = null;
      }
      try {
        tc.custom_params = JSON.parse(this.form.trackingCustomParamsStr || '{}');
      } catch (e) {
        this.$utils.toast(`${this.$t('subscribers.invalidJSON')}: ${e.toString()}`, 'is-danger');
        throw e;
      }
      return tc;
    },
    formatDateTime(s) {
      return dayjs(s).format('YYYY-MM-DD HH:mm');
    },

    onToggleArchivePreview() {
      this.isPreviewingArchive = !this.isPreviewingArchive;
    },

    onAddAltBody() {
      this.form.altbody = htmlToPlainText(this.form.content.body);
    },

    onRemoveAltBody() {
      this.form.altbody = null;
    },

    onShowHeaders() {
      this.isHeadersVisible = !this.isHeadersVisible;
    },

    onShowAttachField() {
      this.isAttachFieldVisible = true;
      this.$nextTick(() => {
        this.$refs.media.focus();
      });
    },

    onOpenAttach() {
      this.isAttachModalOpen = true;
    },

    onAttachSelect(o) {
      if (this.form.media.some((m) => m.id === o.id)) {
        return;
      }

      this.form.media.push(o);
    },

    isUnsaved() {
      return this.data.body !== this.form.content.body
        || this.data.contentType !== this.form.content.contentType;
    },

    onTab(tab) {
      if (tab === 'content' && window.tinymce && window.tinymce.editors.length > 0) {
        this.$nextTick(() => {
          window.tinymce.editors[0].focus();
        });
      }

      // this.$router.replace({ hash: `#${tab}` });
      window.history.replaceState({}, '', `#${tab}`);
    },

    onFillArchiveMeta() {
      const archiveStr = `{"email": "email@domain.com", "name": "${this.$t('globals.fields.name')}", "attribs": {}}`;
      this.form.archiveMetaStr = this.$utils.getPref('campaign.archiveMetaStr') || JSON.stringify(JSON.parse(archiveStr), null, 4);
    },

    onSubmit(typ) {
      // Validate custom JSON headers.
      if (this.form.headersStr && this.form.headersStr !== '[]') {
        try {
          this.form.headers = JSON.parse(this.form.headersStr);
        } catch (e) {
          this.$utils.toast(e.toString(), 'is-danger');
          return undefined;
        }
      } else {
        this.form.headers = [];
      }

      // Validate archive JSON body.
      if (this.form.archive && this.form.archiveMetaStr) {
        try {
          this.form.archiveMeta = JSON.parse(this.form.archiveMetaStr);
        } catch (e) {
          this.$utils.toast(e.toString(), 'is-danger');
          return undefined;
        }
      }

      // Validate custom JSON attribs.
      let attribs = null;
      if (this.form.attribsStr && this.form.attribsStr.trim()) {
        try {
          attribs = JSON.parse(this.form.attribsStr);
        } catch (e) {
          this.$utils.toast(
            `${this.$t('subscribers.invalidJSON')}: ${e.toString()}`,
            'is-danger',

            3000,
          );
          return undefined;
        }
      }
      if (!attribs || typeof attribs !== 'object') attribs = {};
      if (this.form.previewText) attribs.previewText = this.form.previewText;
      else delete attribs.previewText;
      this.form.attribs = attribs;
      this.form.attribsStr = JSON.stringify(attribs, null, 4);

      let trackingConfig = null;
      try {
        trackingConfig = this.buildTrackingConfig();
      } catch {
        return undefined;
      }
      this.form.tracking_config = trackingConfig;

      switch (typ) {
        case 'create':
          return this.createCampaign();
        case 'test':
          return this.sendTest();
        default:
          return this.updateCampaign();
      }
    },

    getCampaign(id) {
      return this.$api.getCampaign(id).then((data) => {
        this.data = data;
        const trackingFields = this.parseTrackingConfig(data.trackingConfig || data.tracking_config || {});
        const { trackingConfig, tracking_config: trackingConfigSnake, ...campaignData } = data;
        this.form = {
          ...this.form,
          ...campaignData,
          ...trackingFields,
          headersStr: JSON.stringify(data.headers, null, 4),
          archiveMetaStr: data.archiveMeta ? JSON.stringify(data.archiveMeta, null, 4) : '{}',
          attribsStr: data.attribs ? JSON.stringify(data.attribs, null, 4) : '{}',

          // The structure that is populated by editor input event.
          content: {
            contentType: data.contentType,
            body: data.body,
            bodySource: data.bodySource,
            templateId: data.templateId,
          },
        };
        this.form.previewText = (data.attribs && data.attribs.previewText) || '';
        this.isAttachFieldVisible = this.form.media.length > 0;

        this.form.media = this.form.media.map((f) => {
          if (!f.id) {
            return { ...f, filename: `❌ ${f.filename}` };
          }
          return f;
        });
      });
    },

    sendTest() {
      let trackingConfig = null;
      try {
        trackingConfig = this.buildTrackingConfig();
      } catch {
        return false;
      }

      const data = {
        id: this.data.id,
        name: this.form.name,
        subject: this.form.subject,
        lists: this.form.lists.map((l) => l.id),
        from_email: this.form.fromEmail,
        reply_to: this.form.replyTo,
        messenger: this.form.messenger,
        type: 'regular',
        headers: this.form.headers,
        tags: this.form.tags,
        template_id: this.form.content.templateId,
        content_type: this.form.content.contentType,
        body: this.form.content.body,
        altbody: this.form.content.contentType !== 'plain' ? this.form.altbody : null,
        subscribers: this.form.testEmails,
        media: this.form.media.map((m) => m.id),
        tracking_config: trackingConfig,
      };

      this.$api.testCampaign(data).then(() => {
        this.$utils.toast(this.$t('campaigns.testSent'));
      });
      return false;
    },

    createCampaign() {
      const data = {
        archiveSlug: this.form.subject,
        name: this.form.name,
        subject: this.form.subject,
        lists: this.form.lists.map((l) => l.id),
        from_email: this.form.fromEmail,
        reply_to: this.form.replyTo,
        content_type: this.form.content.contentType,
        messenger: this.form.messenger,
        type: 'regular',
        tags: this.form.tags,
        send_at: this.form.sendLater ? this.form.sendAtDate : null,
        headers: this.form.headers,
        attribs: this.form.attribs,
        tracking_config: this.form.tracking_config,
        media: this.form.media.map((m) => m.id),
      };

      return this.$api.createCampaign(data).then((d) => {
        this.isNew = false;
        this.isEditing = true;
        this.$router.replace({ name: 'campaign', params: { id: d.id } });
        return this.getCampaign(d.id);
      });
    },

    async updateCampaign(typ) {
      try {
        this.form.tracking_config = this.buildTrackingConfig();
      } catch {
        return Promise.reject();
      }

      const data = {
        archive_slug: this.form.archiveSlug,
        name: this.form.name,
        subject: this.form.subject,
        lists: this.form.lists.map((l) => l.id),
        from_email: this.form.fromEmail,
        reply_to: this.form.replyTo,
        messenger: this.form.messenger,
        type: 'regular',
        tags: this.form.tags,
        send_at: this.form.sendLater ? this.form.sendAtDate : null,
        headers: this.form.headers,
        attribs: this.form.attribs,
        template_id: this.form.content.templateId,
        content_type: this.form.content.contentType,
        body: this.form.content.body,
        body_source: this.form.content.bodySource,
        altbody: this.form.content.contentType !== 'plain' ? this.form.altbody : null,
        archive: this.form.archive,
        archive_template_id: this.form.archiveTemplateId,
        archive_meta: this.form.archiveMeta,
        tracking_config: this.form.tracking_config,
        media: this.form.media.map((m) => m.id),
      };

      let typMsg = 'globals.messages.updated';
      if (typ === 'start') {
        typMsg = 'campaigns.started';
      }

      if (!this.form.sendAtDate) {
        this.form.sendLater = false;
      }

      // This promise is used by startCampaign to first save before starting.
      return new Promise((resolve) => {
        this.$api.updateCampaign(this.data.id, data).then((d) => {
          this.data = d;
          this.form.archiveSlug = d.archiveSlug;
          this.form.attribsStr = d.attribs ? JSON.stringify(d.attribs, null, 4) : '{}';

          this.$utils.toast(this.$t(typMsg, { name: d.name }));
          resolve();
        });
      });
    },

    onUpdateCampaignArchive() {
      if (this.isEditing && this.canEdit) {
        return;
      }

      const data = {
        archive: this.form.archive,
        archive_template_id: this.form.archiveTemplateId,
        archive_meta: JSON.parse(this.form.archiveMetaStr),
        archive_slug: this.form.archiveSlug,
      };

      this.$api.updateCampaignArchive(this.data.id, data).then((d) => {
        this.form.archiveSlug = d.archiveSlug;
      });
    },

    // Starts or schedule a campaign.
    startCampaign(opts) {
      if (!this.canStart && !this.canSchedule) {
        return;
      }

      const run = () => {
        this.updateCampaign().then(() => {
          let status = '';
          if (this.canStart) {
            status = 'running';
          } else if (this.canSchedule) {
            status = 'scheduled';
          } else {
            return;
          }

          this.$api.changeCampaignStatus(this.data.id, status).then(() => {
            this.$router.push({ name: 'campaigns' });
          });
        });
      };

      if (opts && opts.skipConfirm) {
        run();
        return;
      }

      this.$utils.confirm(null, run);
    },

    unscheduleCampaign() {
      this.$api.changeCampaignStatus(this.data.id, 'draft').then((d) => {
        this.data = d;
      });
    },

    applyCreateQuery() {
      const q = this.$route.query || {};
      if (q.name) this.form.name = String(q.name);
      if (q.subject) this.form.subject = String(q.subject);
      if (q.messenger) {
        this.form.messenger = String(q.messenger);
      } else if (q.channel && q.channel !== 'email') {
        const ch = String(q.channel).toLowerCase();
        const messengers = (this.serverConfig && this.serverConfig.messengers) || [];
        const found = messengers.find((m) => String(m).toLowerCase().indexOf(ch) > -1);
        if (found) this.form.messenger = found;
      }
      if (q.tags) {
        this.form.tags = String(q.tags).split(',').map((t) => t.trim()).filter(Boolean);
      }
    },

    showSetupMissing() {
      if (!this.form.name || !this.form.name.trim()) {
        this.form.name = 'Untitled campaign';
      }
      if (!this.form.subject || !this.form.subject.trim()) {
        this.$utils.toast('Add a subject before continuing.');
        return true;
      }
      if (!this.form.lists || !this.form.lists.length) {
        this.$utils.toast('Select recipients before continuing.');
        return true;
      }
      return false;
    },

    ensureSaved() {
      if (this.isNew) {
        if (this.showSetupMissing()) return Promise.reject();
        const created = this.onSubmit('create');
        if (!created) return Promise.reject();
        return created;
      }
      return this.updateCampaign();
    },

    onSetupSave(opts) {
      if (this.isNew && !(opts && opts.createIfNew)) return;
      if (this.isNew) {
        if (this.showSetupMissing()) return;
        this.onSubmit('create');
        return;
      }
      this.onSubmit('update');
    },

    onEnsureSaved(cb) {
      this.ensureSaved().then(() => {
        if (typeof cb === 'function') cb();
      }).catch(() => {});
    },

    onSetupSchedule(opts) {
      if (opts && opts.sendNow) {
        this.form.sendLater = false;
        this.form.sendAtDate = null;
      } else {
        this.form.sendLater = true;
        this.form.sendAtDate = opts && opts.at ? opts.at : this.form.sendAtDate;
      }
      this.ensureSaved().then(() => {
        this.startCampaign({ skipConfirm: true });
      }).catch(() => {});
    },

    onUnschedule() {
      this.$utils.confirm(null, this.unscheduleCampaign);
    },

    goToEditor() {
      this.$router.replace({
        name: 'campaign',
        params: { id: this.data.id },
        query: { view: 'edit' },
        hash: this.$route.hash || '',
      });
    },
  },

  computed: {
    ...mapState(['serverConfig', 'loading', 'lists', 'templates']),

    newCampaignTitle() {
      const q = this.$route.query || {};
      if (q.template === 'automation' || q.template) {
        return q.name || 'New automation';
      }
      if (q.channel === 'sms') return 'New SMS campaign';
      if (q.channel === 'whatsapp') return 'New WhatsApp campaign';
      if (q.channel === 'push') return 'New push campaign';
      return this.$t('campaigns.newCampaign');
    },

    newCampaignLead() {
      const q = this.$route.query || {};
      const leads = {
        welcome: 'Welcome automation — send a message when someone joins a list.',
        abandoned_cart: 'Abandoned cart — follow up after a contact leaves items behind.',
        marketing_activity: 'Marketing activity — follow up when contacts open or click.',
        product_purchase: 'Product purchase — send a message after a purchase.',
        anniversary: 'Anniversary — send messages for a special date or birthday.',
        automation: 'Create an automation from scratch.',
      };
      if (q.template && leads[q.template]) return leads[q.template];
      if (q.channel === 'sms') return 'Create a new SMS campaign.';
      if (q.channel === 'whatsapp') return 'Create a new WhatsApp campaign.';
      if (q.channel === 'push') return 'Create a new push notification campaign.';
      return 'Create a new email campaign for your lists.';
    },

    showSetup() {
      return !this.showReport;
    },

    allLists() {
      return (this.lists && this.lists.results) || [];
    },

    showReport() {
      if (!this.data || !this.data.id) return false;
      if (this.$route.query.view === 'edit') return false;
      if (this.$route.query.view === 'report') return true;
      return ['finished', 'cancelled', 'running'].indexOf(this.data.status) > -1;
    },

    canViewReport() {
      return !!this.data.id;
    },

    canManage() {
      return this.$can('campaigns:manage_all', 'campaigns:manage');
    },

    canSend() {
      return this.$can('campaigns:send');
    },

    canEdit() {
      return this.isNew
        || this.data.status === 'draft' || this.data.status === 'scheduled' || this.data.status === 'paused';
    },

    canSchedule() {
      return (this.data.status === 'draft' || this.data.status === 'paused') && (this.form.sendLater && this.form.sendAtDate);
    },

    canUnSchedule() {
      return this.data.status === 'scheduled';
    },

    canStart() {
      return (this.data.status === 'draft' || this.data.status === 'paused') && !this.form.sendLater;
    },

    canDownloadReport() {
      return ['finished', 'paused', 'running'].includes(this.data.status) && this.data.sent > 0;
    },

    canArchive() {
      return this.data.status !== 'cancelled' && this.data.type !== 'optin';
    },

    selectedLists() {
      if (this.selListIDs.length === 0 || !this.lists.results) {
        return [];
      }

      return this.lists.results.filter((l) => this.selListIDs.indexOf(l.id) > -1);
    },

    emailMessengers() {
      return ['email', ...this.serverConfig.messengers.filter((m) => m.startsWith('email-'))];
    },

    otherMessengers() {
      return this.serverConfig.messengers.filter((m) => m !== 'email' && !m.startsWith('email-'));
    },
  },

  beforeRouteLeave(to, from, next) {
    if (this.isUnsaved()) {
      this.$utils.confirm(this.$t('globals.messages.confirmDiscard'), () => next(true));
      return;
    }
    next(true);
  },

  watch: {
    selectedLists() {
      this.form.lists = this.selectedLists;
    },

    // eslint-disable-next-line func-names
    'data.sendAt': function () {
      if (this.data.sendAt !== null) {
        this.form.sendLater = true;
        this.form.sendAtDate = dayjs(this.data.sendAt).toDate();
      } else {
        this.form.sendLater = false;
        this.form.sendAtDate = null;
      }
    },
  },

  mounted() {
    window.onbeforeunload = () => this.isUnsaved() || null;

    // Fill default form fields.
    this.form.fromEmail = this.serverConfig.from_email;

    // New campaign.
    const { id } = this.$route.params;
    if (id === 'new') {
      this.isNew = true;

      if (this.$route.query.list_id) {
        // Multiple list_id query params.
        let strIds = [];
        if (typeof this.$route.query.list_id === 'object') {
          strIds = this.$route.query.list_id;
        } else {
          strIds = [this.$route.query.list_id];
        }

        this.selListIDs = strIds.map((v) => parseInt(v, 10));
      }
    } else {
      const intID = parseInt(id, 10);
      if (intID <= 0 || Number.isNaN(intID)) {
        this.$utils.toast(this.$t('campaigns.invalid'));
        return;
      }

      this.isEditing = true;
    }

    this.$api.getLists({ per_page: 'all', status: 'active' }).catch(() => {});

    // Get templates list.
    this.$api.getTemplates().then((data) => {
      if (data.length > 0) {
        if (!this.form.templateId) {
          const tpl = data.find((i) => i.isDefault === true);
          this.form.templateId = tpl.id;
        }
      }
    });

    // Fetch campaign.
    if (this.isEditing) {
      this.getCampaign(id).then(() => {
        if (this.$route.hash !== '' && this.$route.query.view !== 'report') {
          this.activeTab = this.$route.hash.replace('#', '');
        }
      });
    } else {
      this.form.messenger = 'email';
      this.applyCreateQuery();
    }

    this.$nextTick(() => {
      if (this.$refs.focus && this.$refs.focus.focus) {
        this.$refs.focus.focus();
      }
    });

    this.$events.$on('campaign.update', () => {
      this.onSubmit('update');
    });
  },

  beforeDestroy() {
    this.$events.$off('campaign.update');
  },
});
</script>
