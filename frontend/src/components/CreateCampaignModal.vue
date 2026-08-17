<template>
  <div class="create-campaign-modal" data-cy="create-campaign-modal">
    <header class="create-campaign-modal__head">
      <h2>Create a campaign</h2>
      <button type="button" class="create-campaign-modal__close" aria-label="Close" @click="$emit('close')">
        <svg width="14" height="14" viewBox="0 0 14 14" fill="none" aria-hidden="true">
          <path d="M2 2l10 10M12 2L2 12" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
        </svg>
      </button>
    </header>

    <div class="create-campaign-modal__body">
      <section class="create-campaign-modal__section create-campaign-modal__section--standard">
        <h3>Standard</h3>
        <p>Create a one-off campaign from scratch.</p>
        <div class="create-campaign-modal__channels">
          <button
            v-for="ch in standardChannels"
            :key="ch.id"
            type="button"
            class="create-campaign-modal__channel"
            :class="'is-' + ch.id"
            @click="onStandard(ch.id)"
          >
            <span class="create-campaign-modal__channel-icon" aria-hidden="true">
              <svg v-if="ch.id === 'email'" width="54" height="42" viewBox="0 0 54 42" fill="none">
                <rect x="26" y="2" width="22" height="26" rx="2.5" fill="#fff" stroke="#0B996E" stroke-width="1.7" />
                <path d="M30.5 10h13M30.5 15h9" stroke="#0B996E" stroke-width="1.6" stroke-linecap="round" />
                <path d="M4 12h32v24H4V12z" fill="#E9F8EE" stroke="#0B996E" stroke-width="1.8" />
                <path d="M4 12l16 11.5L36 12" fill="#fff" stroke="#0B996E" stroke-width="1.8" stroke-linejoin="round" />
              </svg>
              <svg v-else-if="ch.id === 'sms'" width="48" height="42" viewBox="0 0 48 42" fill="none">
                <path d="M8 6h28a8 8 0 018 8v8a8 8 0 01-8 8H18l-10 8v-8H8a8 8 0 01-8-8v-8a8 8 0 018-8z"
                  fill="#7BC9A6" />
                <circle cx="16" cy="20" r="2.2" fill="#fff" />
                <circle cx="24" cy="20" r="2.2" fill="#fff" />
                <circle cx="32" cy="20" r="2.2" fill="#fff" />
              </svg>
              <svg v-else-if="ch.id === 'whatsapp'" width="42" height="42" viewBox="0 0 42 42" fill="none">
                <circle cx="21" cy="21" r="19" fill="#25D366" />
                <path
                  fill="#fff"
                  d="M13.2 29.4l.8-3A9 9 0 1124.8 30l3.2.9-3.5-.6a9 9 0 01-11.3-.9z"
                />
                <path
                  fill="#25D366"
                  d="M21 14.2c-3.9 0-7 3.1-7 7 0 1.3.4 2.5 1 3.5l-.7 2.6 2.7-.7c1 .6 2.1.9 3.3.9 3.9 0 7-3.1 7-7s-3.1-7.3-7.3-7.3z"
                />
              </svg>
              <svg v-else-if="ch.id === 'push'" width="52" height="42" viewBox="0 0 52 42" fill="none">
                <rect x="2" y="10" width="32" height="24" rx="3" fill="#fff" stroke="#C9B48A" stroke-width="1.6" />
                <path d="M2 16h32" stroke="#C9B48A" stroke-width="1.5" />
                <circle cx="7" cy="13" r="1.1" fill="#C9B48A" />
                <circle cx="11" cy="13" r="1.1" fill="#C9B48A" />
                <rect x="26" y="3" width="24" height="18" rx="4" fill="#E9F8EE" stroke="#0B996E" stroke-width="1.6" />
                <path d="M38 8.2a3.4 3.4 0 013.4 3.4c0 2.1.9 3.2.9 3.2h-8.6s.9-1.1.9-3.2A3.4 3.4 0 0138 8.2z"
                  fill="#fff" stroke="#0B996E" stroke-width="1.4" />
                <circle cx="38" cy="17.6" r="1" fill="#0B996E" />
              </svg>
              <svg v-else width="52" height="42" viewBox="0 0 52 42" fill="none">
                <rect x="2" y="6" width="34" height="26" rx="3" fill="#fff" stroke="#8BB8D6" stroke-width="1.6" />
                <path d="M2 12h34" stroke="#8BB8D6" stroke-width="1.5" />
                <rect x="16" y="14" width="34" height="22" rx="3.5" fill="#E9F8EE" stroke="#0B996E" stroke-width="1.7" />
                <path d="M44 20h4M46 18v4" stroke="#0B996E" stroke-width="1.6" stroke-linecap="round" />
                <path d="M22 24h20M22 29h14" stroke="#0B996E" stroke-width="1.5" stroke-linecap="round" />
              </svg>
            </span>
            <span class="create-campaign-modal__channel-label">
              {{ ch.label }}
              <span v-if="ch.badge === 'activate'" class="create-campaign-modal__badge">Activate</span>
              <svg
                v-else-if="ch.badge === 'crown'"
                class="create-campaign-modal__crown"
                width="14"
                height="12"
                viewBox="0 0 14 12"
                aria-hidden="true"
              >
                <path d="M1.2 10.4h11.6L11.6 4.2 9 6.6 7 1.6 5 6.6 2.4 4.2z" fill="#E8B923" />
                <rect x="1.1" y="10.1" width="11.8" height="1.5" rx="0.5" fill="#D4A017" />
              </svg>
            </span>
          </button>
        </div>
      </section>

      <section class="create-campaign-modal__section create-campaign-modal__section--automated" ref="automated">
        <div class="create-campaign-modal__section-head">
          <div>
            <h3>Automated</h3>
            <p>
              Create an automation from scratch, get assistance from our AI,
              or pick one of our pre-built automations.
            </p>
          </div>
          <button type="button" class="create-campaign-modal__scratch" @click="onScratch">
            Create from scratch
          </button>
        </div>

        <button
          v-if="!aiOpen"
          type="button"
          class="create-campaign-modal__ai-teaser"
          @click="aiOpen = true"
        >
          Start with a sentence, we'll build your automation
        </button>

        <template v-else>
          <div class="create-campaign-modal__ai-banner">
            <label class="is-sr-only" for="ai-write">Describe your automation</label>
            <textarea
              id="ai-write"
              v-model="aiText"
              class="create-campaign-modal__ai-banner-input"
              rows="3"
              placeholder="Start with a sentence, we'll build your automation"
              @keydown.enter.ctrl.prevent="submitBanner"
            />
            <div class="create-campaign-modal__ai-banner-actions">
              <button type="button" class="create-campaign-modal__ai-circle" aria-label="Reset" @click="resetAi">
                <svg width="18" height="18" viewBox="0 0 18 18" fill="none" aria-hidden="true">
                  <path d="M4 9a5 5 0 019.6-2M14 9a5 5 0 01-9.6 2" stroke="#6B7280" stroke-width="1.6" stroke-linecap="round" />
                  <path d="M13.4 3.2V6.4H10.2M4.6 14.8v-3.2H7.8" stroke="#6B7280" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" />
                </svg>
              </button>
              <button type="button" class="create-campaign-modal__ai-circle" aria-label="Create automation" @click="submitBanner">
                <svg width="18" height="18" viewBox="0 0 18 18" fill="none" aria-hidden="true">
                  <path d="M9 13.5V4.5M5.2 8.2L9 4.5l3.8 3.7" stroke="#6B7280" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
                </svg>
              </button>
            </div>
          </div>

          <div class="create-campaign-modal__filters">
            <div class="create-campaign-modal__filter">
              <span id="prebuilt-type-label">Pre-built automation type</span>
              <span class="create-campaign-modal__token-wrap is-filter">
                <button
                  type="button"
                  class="create-campaign-modal__select"
                  aria-labelledby="prebuilt-type-label"
                  @click.stop="toggleMenu('typeFilter')"
                >
                  {{ typeFilterLabel }}
                  <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
                    <path d="M3 4.2L6 7.2 9 4.2" stroke="currentColor" stroke-width="1.5"
                      stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                </button>
                <ul v-if="openMenu === 'typeFilter'" class="create-campaign-modal__menu">
                  <li v-for="opt in typeFilterOptions" :key="opt.value">
                    <button type="button" @click="typeFilter = opt.value; openMenu = null">{{ opt.label }}</button>
                  </li>
                </ul>
              </span>
            </div>
            <div class="create-campaign-modal__filter">
              <span id="prebuilt-channel-label">Channel</span>
              <span class="create-campaign-modal__token-wrap is-filter">
                <button
                  type="button"
                  class="create-campaign-modal__select"
                  aria-labelledby="prebuilt-channel-label"
                  @click.stop="toggleMenu('channelFilter')"
                >
                  {{ channelFilterLabel }}
                  <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
                    <path d="M3 4.2L6 7.2 9 4.2" stroke="currentColor" stroke-width="1.5"
                      stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                </button>
                <ul v-if="openMenu === 'channelFilter'" class="create-campaign-modal__menu">
                  <li v-for="opt in channelFilterOptions" :key="opt.value">
                    <button type="button" @click="channelFilter = opt.value; openMenu = null">{{ opt.label }}</button>
                  </li>
                </ul>
              </span>
            </div>
          </div>

          <div class="create-campaign-modal__templates">
            <button
              v-for="tpl in visibleTemplates"
              :key="tpl.id"
              type="button"
              class="create-campaign-modal__template"
              @click="onTemplate(tpl)"
            >
              <strong>{{ tpl.title }}</strong>
              <span>{{ tpl.desc }}</span>
            </button>
            <p v-if="!visibleTemplates.length" class="create-campaign-modal__empty">
              No pre-built automations match these filters.
            </p>
          </div>
        </template>
      </section>
    </div>
  </div>
</template>

<script>
import Vue from 'vue';
import { mapState } from 'vuex';

const TEMPLATES = [
  {
    id: 'welcome',
    title: 'Welcome message',
    desc: 'Send a welcome message after a subscriber joins your list.',
    name: 'Welcome message',
    subject: 'Welcome — thanks for joining',
    tags: 'welcome,automation',
    channels: ['email', 'sms', 'whatsapp'],
  },
  {
    id: 'abandoned_cart',
    title: 'Abandoned cart',
    desc: 'Send a message after a contact abandons a cart.',
    name: 'Abandoned cart',
    subject: 'You left something behind',
    tags: 'abandoned-cart,automation',
    channels: ['email', 'sms', 'whatsapp'],
  },
  {
    id: 'marketing_activity',
    title: 'Marketing activity',
    desc: 'Send messages based on whether contacts open or click an email campaign.',
    name: 'Marketing activity',
    subject: 'Following up on our last email',
    tags: 'engagement,automation',
    channels: ['email', 'sms', 'whatsapp'],
  },
  {
    id: 'product_purchase',
    title: 'Product purchase',
    desc: 'Send a message when a product is purchased on your website.',
    name: 'Product purchase',
    subject: 'Thanks for your purchase',
    tags: 'purchase,automation',
    channels: ['email', 'sms', 'whatsapp'],
  },
  {
    id: 'anniversary',
    title: 'Anniversary date',
    desc: 'Send a series of messages based on a special event or birthday.',
    name: 'Anniversary date',
    subject: 'Happy anniversary',
    tags: 'anniversary,automation',
    channels: ['email', 'sms', 'whatsapp'],
  },
];

export default Vue.extend({
  name: 'CreateCampaignModal',

  props: {
    startSection: {
      type: String,
      default: 'standard',
    },
  },

  data() {
    return {
      aiOpen: false,
      aiMode: 'write',
      aiAction: 'sends',
      aiChannel: '',
      aiTrigger: '',
      aiText: '',
      openMenu: null,
      typeFilter: 'popular',
      channelFilter: 'email',
      actionOptions: [
        { value: 'sends', label: 'sends' },
        { value: 'does_not_send', label: "doesn't send" },
      ],
      channelOptions: [
        { value: 'email', label: 'email' },
        { value: 'sms', label: 'SMS' },
        { value: 'whatsapp', label: 'WhatsApp' },
      ],
      triggerOptions: [
        { value: 'welcome', label: 'a subscriber joins a list' },
        { value: 'abandoned_cart', label: 'a contact abandons a cart' },
        { value: 'marketing_activity', label: 'a contact opens or clicks an email' },
        { value: 'product_purchase', label: 'a product is purchased' },
        { value: 'anniversary', label: 'an anniversary or birthday arrives' },
      ],
      typeFilterOptions: [
        { value: 'popular', label: 'Most popular' },
        { value: 'welcome', label: 'Welcome message' },
        { value: 'abandoned_cart', label: 'Abandoned cart' },
        { value: 'marketing_activity', label: 'Marketing activity' },
        { value: 'product_purchase', label: 'Product purchase' },
        { value: 'anniversary', label: 'Anniversary date' },
      ],
      channelFilterOptions: [
        { value: 'email', label: 'Email' },
        { value: 'sms', label: 'SMS' },
        { value: 'whatsapp', label: 'WhatsApp' },
      ],
      standardChannels: [
        { id: 'email', label: 'Email', badge: '' },
        { id: 'sms', label: 'SMS', badge: '' },
        { id: 'whatsapp', label: 'WhatsApp', badge: 'crown' },
        { id: 'push', label: 'Push', badge: 'activate' },
        { id: 'popup', label: 'Pop-up', badge: 'crown' },
      ],
      templates: TEMPLATES,
    };
  },

  computed: {
    ...mapState(['serverConfig']),

    actionLabel() {
      const found = this.actionOptions.find((o) => o.value === this.aiAction);
      return found ? found.label : 'sends';
    },

    channelTokenLabel() {
      if (!this.aiChannel) return 'channel';
      const found = this.channelOptions.find((o) => o.value === this.aiChannel);
      return found ? found.label : 'channel';
    },

    triggerTokenLabel() {
      if (!this.aiTrigger) return 'something happens';
      const found = this.triggerOptions.find((o) => o.value === this.aiTrigger);
      return found ? found.label : 'something happens';
    },

    typeFilterLabel() {
      const found = this.typeFilterOptions.find((o) => o.value === this.typeFilter);
      return found ? found.label : 'Most popular';
    },

    channelFilterLabel() {
      const found = this.channelFilterOptions.find((o) => o.value === this.channelFilter);
      return found ? found.label : 'Email';
    },

    visibleTemplates() {
      return this.templates.filter((tpl) => {
        const typeOk = this.typeFilter === 'popular' || tpl.id === this.typeFilter;
        const chOk = !this.channelFilter || tpl.channels.indexOf(this.channelFilter) > -1;
        return typeOk && chOk;
      });
    },
  },

  mounted() {
    document.addEventListener('click', this.onDocClick);
    if (this.startSection === 'automated') {
      this.$nextTick(() => {
        if (this.$refs.automated && this.$refs.automated.scrollIntoView) {
          this.$refs.automated.scrollIntoView({ block: 'start' });
        }
      });
    }
  },

  beforeDestroy() {
    document.removeEventListener('click', this.onDocClick);
  },

  methods: {
    onDocClick() {
      this.openMenu = null;
    },

    toggleMenu(name) {
      this.openMenu = this.openMenu === name ? null : name;
    },

    resetAi() {
      this.aiAction = 'sends';
      this.aiChannel = '';
      this.aiTrigger = '';
      this.aiText = '';
      this.openMenu = null;
      this.typeFilter = 'popular';
      this.channelFilter = 'email';
    },

    submitBanner() {
      this.aiMode = 'write';
      if (!(this.aiText || '').trim()) {
        this.onScratch();
        return;
      }
      this.submitAi();
    },

    messengers() {
      const cfg = this.serverConfig || {};
      return cfg.messengers || [];
    },

    resolveMessenger(channel) {
      if (!channel || channel === 'email') return 'email';
      const list = this.messengers();
      const ch = String(channel).toLowerCase();
      const keys = {
        sms: ['sms', 'text', 'twilio', 'msg91', 'vonage', 'nexmo'],
        whatsapp: ['whatsapp', 'wa-'],
        push: ['push', 'fcm', 'onesignal', 'webpush'],
      };
      const needles = keys[ch] || [ch];
      const found = list.find((m) => {
        const n = String(m).toLowerCase();
        return needles.some((k) => n.indexOf(k) > -1);
      });
      return found || null;
    },

    goToCampaign(query) {
      this.$emit('close');
      this.$router.push({
        name: 'campaign',
        params: { id: 'new' },
        query,
      });
    },

    onStandard(id) {
      if (id === 'popup') {
        this.$emit('close');
        this.$router.push({ name: 'forms' });
        return;
      }

      if (id === 'push') {
        const messenger = this.resolveMessenger('push');
        if (!messenger) {
          this.$utils.toast('Activate Push by adding a push messenger in Settings.');
          this.$emit('close');
          this.$router.push({ name: 'settings' });
          return;
        }
        this.goToCampaign({
          channel: 'push',
          messenger,
          name: 'Push campaign',
          subject: 'New notification',
        });
        return;
      }

      if (id === 'sms' || id === 'whatsapp') {
        const messenger = this.resolveMessenger(id);
        if (!messenger) {
          this.$utils.toast(
            `No ${id === 'sms' ? 'SMS' : 'WhatsApp'} messenger configured. `
            + 'Pick a messenger in the campaign or add one in Settings.',
          );
        }
        const query = {
          channel: id,
          name: id === 'sms' ? 'SMS campaign' : 'WhatsApp campaign',
          subject: id === 'sms' ? 'New SMS' : 'New WhatsApp message',
        };
        if (messenger) query.messenger = messenger;
        this.goToCampaign(query);
        return;
      }

      this.goToCampaign({ channel: 'email' });
    },

    onScratch() {
      this.goToCampaign({
        channel: 'email',
        template: 'automation',
        name: 'New automation',
        subject: 'New automation',
        tags: 'automation',
      });
    },

    templateQuery(tpl, channel) {
      const messenger = this.resolveMessenger(channel);
      const query = {
        channel: channel || 'email',
        template: tpl.id,
        name: tpl.name,
        subject: tpl.subject,
        tags: tpl.tags,
      };
      if (messenger && messenger !== 'email') query.messenger = messenger;
      return query;
    },

    onTemplate(tpl) {
      const channel = this.channelFilter || 'email';
      if (channel !== 'email' && !this.resolveMessenger(channel)) {
        this.$utils.toast(
          `No ${channel === 'sms' ? 'SMS' : 'WhatsApp'} messenger configured. `
          + 'The campaign will open with email — change messenger if needed.',
        );
      }
      this.goToCampaign(this.templateQuery(tpl, channel));
    },

    matchWriteTemplate(text) {
      const lower = text.toLowerCase();
      const byTitle = this.templates.find((t) => lower.indexOf(t.title.toLowerCase()) > -1);
      if (byTitle) return byTitle;
      if (lower.indexOf('welcome') > -1) {
        return this.templates.find((t) => t.id === 'welcome');
      }
      if (lower.indexOf('cart') > -1) {
        return this.templates.find((t) => t.id === 'abandoned_cart');
      }
      if (lower.indexOf('purchase') > -1) {
        return this.templates.find((t) => t.id === 'product_purchase');
      }
      if (lower.indexOf('birthday') > -1 || lower.indexOf('anniversary') > -1) {
        return this.templates.find((t) => t.id === 'anniversary');
      }
      return {
        id: 'custom',
        name: text.slice(0, 80),
        subject: text.slice(0, 120),
        tags: 'automation',
      };
    },

    submitAi() {
      if (this.aiMode === 'write') {
        const text = (this.aiText || '').trim();
        if (!text) {
          this.$utils.toast('Describe the automation you want to create.');
          return;
        }
        const lower = text.toLowerCase();
        let channel = 'email';
        if (lower.indexOf('whatsapp') > -1) {
          channel = 'whatsapp';
        } else if (lower.indexOf('sms') > -1 || lower.indexOf('text message') > -1) {
          channel = 'sms';
        }
        this.goToCampaign(this.templateQuery(this.matchWriteTemplate(text), channel));
        return;
      }

      if (!this.aiChannel || !this.aiTrigger) {
        this.$utils.toast('Choose a channel and when something happens.');
        return;
      }

      const tpl = this.templates.find((t) => t.id === this.aiTrigger) || {
        id: this.aiTrigger || 'custom',
        name: 'New automation',
        subject: 'New automation',
        tags: 'automation',
      };
      const query = this.templateQuery(tpl, this.aiChannel);
      if (this.aiAction === 'does_not_send') {
        query.name = `${tpl.name} (hold)`;
        query.subject = tpl.subject;
      }
      this.goToCampaign(query);
    },
  },
});
</script>
