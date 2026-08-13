<template>
  <div class="create-campaign-modal">
    <header class="create-campaign-modal__head">
      <h2>Create a campaign</h2>
      <button type="button" class="create-campaign-modal__close" aria-label="Close" @click="$emit('close')">
        <svg width="14" height="14" viewBox="0 0 14 14" fill="none" aria-hidden="true">
          <path d="M2 2l10 10M12 2L2 12" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
        </svg>
      </button>
    </header>

    <div class="create-campaign-modal__body">
      <section class="create-campaign-modal__section">
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
              <svg v-if="ch.id === 'email'" width="36" height="28" viewBox="0 0 36 28" fill="none">
                <rect x="4" y="6" width="24" height="16" rx="2" fill="#fff" stroke="#7BC9A6" stroke-width="1.6" />
                <path d="M4 8l12 8 12-8" stroke="#7BC9A6" stroke-width="1.6" fill="none" />
                <rect x="18" y="12" width="14" height="12" rx="1.5" fill="#E8F7EC" stroke="#0B996E" stroke-width="1.4" />
                <path d="M21 17h8M21 20h5" stroke="#0B996E" stroke-width="1.3" stroke-linecap="round" />
              </svg>
              <svg v-else-if="ch.id === 'sms'" width="32" height="28" viewBox="0 0 32 28" fill="none">
                <rect x="3" y="4" width="26" height="16" rx="8" fill="#fff" stroke="#9B8CFF" stroke-width="1.6" />
                <circle cx="11" cy="12" r="1.6" fill="#7BC9A6" />
                <circle cx="16" cy="12" r="1.6" fill="#7BC9A6" />
                <circle cx="21" cy="12" r="1.6" fill="#7BC9A6" />
                <path d="M10 20l-4 5 7-3" fill="#fff" stroke="#9B8CFF" stroke-width="1.4" stroke-linejoin="round" />
              </svg>
              <svg v-else-if="ch.id === 'whatsapp'" width="30" height="30" viewBox="0 0 30 30" fill="none">
                <circle cx="15" cy="15" r="12" fill="#25D366" />
                <path d="M9.2 20.6l.9-3.2A7 7 0 1118.8 21l3.1.9-3.4-.6a7 7 0 01-9.3-.7z" fill="#fff" />
              </svg>
              <svg v-else-if="ch.id === 'push'" width="34" height="28" viewBox="0 0 34 28" fill="none">
                <rect x="3" y="6" width="22" height="16" rx="2" fill="#fff" stroke="#D4B483" stroke-width="1.5" />
                <path d="M3 10h22" stroke="#D4B483" stroke-width="1.4" />
                <circle cx="6.5" cy="8" r="0.9" fill="#D4B483" />
                <circle cx="9.2" cy="8" r="0.9" fill="#D4B483" />
                <rect x="18" y="3" width="13" height="10" rx="2" fill="#E8F7EC" stroke="#0B996E" stroke-width="1.4" />
                <path d="M24.5 6.2a2.2 2.2 0 012.2 2.2c0 1.4.6 2.1.6 2.1h-5.6s.6-.7.6-2.1a2.2 2.2 0 012.2-2.2z"
                  fill="#fff" stroke="#0B996E" stroke-width="1.2" />
                <circle cx="24.5" cy="12.2" r="0.8" fill="#0B996E" />
              </svg>
              <svg v-else width="34" height="28" viewBox="0 0 34 28" fill="none">
                <rect x="2" y="5" width="24" height="18" rx="2" fill="#fff" stroke="#7EB6D9" stroke-width="1.5" />
                <rect x="10" y="9" width="20" height="14" rx="2" fill="#E8F7EC" stroke="#0B996E" stroke-width="1.5" />
                <path d="M26.5 12.2h2.2M27.6 11.1v2.2" stroke="#0B996E" stroke-width="1.3" stroke-linecap="round" />
                <path d="M14 15h10M14 18h7" stroke="#0B996E" stroke-width="1.3" stroke-linecap="round" />
              </svg>
            </span>
            <span class="create-campaign-modal__channel-label">
              {{ ch.label }}
              <span v-if="ch.badge === 'activate'" class="create-campaign-modal__badge">Activate</span>
              <span v-else-if="ch.badge === 'crown'" class="create-campaign-modal__crown" aria-hidden="true">♛</span>
            </span>
          </button>
        </div>
      </section>

      <section class="create-campaign-modal__section" ref="automated">
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
          Start with a sentence. we’ll build your automation
        </button>

        <div v-else class="create-campaign-modal__ai">
          <h4>Start with a sentence, we’ll build your automation</h4>
          <div class="create-campaign-modal__ai-modes" role="tablist" aria-label="AI mode">
            <button
              type="button"
              class="create-campaign-modal__ai-mode"
              :class="{ 'is-active': aiMode === 'guided' }"
              @click="aiMode = 'guided'"
            >
              <svg width="14" height="14" viewBox="0 0 14 14" fill="none" aria-hidden="true">
                <path d="M2 7h10M7 2v10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
                <circle cx="7" cy="7" r="5.2" stroke="currentColor" stroke-width="1.4" />
              </svg>
              Guided mode
            </button>
            <button
              type="button"
              class="create-campaign-modal__ai-mode"
              :class="{ 'is-active': aiMode === 'write' }"
              @click="aiMode = 'write'"
            >
              <svg width="14" height="14" viewBox="0 0 14 14" fill="none" aria-hidden="true">
                <path d="M8.2 2.4l3.4 3.4L4.6 13H1.2v-3.4L8.2 2.4z" stroke="currentColor" stroke-width="1.4"
                  stroke-linejoin="round" />
              </svg>
              Write your own
            </button>
          </div>

          <div v-if="aiMode === 'guided'" class="create-campaign-modal__ai-box">
            <p class="create-campaign-modal__sentence">
              Create an automation that
              <span class="create-campaign-modal__token-wrap">
                <button type="button" class="create-campaign-modal__token" @click.stop="toggleMenu('action')">
                  {{ actionLabel }}
                  <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
                    <path d="M3 4.2L6 7.2 9 4.2" stroke="currentColor" stroke-width="1.5"
                      stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                </button>
                <ul v-if="openMenu === 'action'" class="create-campaign-modal__menu">
                  <li v-for="opt in actionOptions" :key="opt.value">
                    <button type="button" @click="aiAction = opt.value; openMenu = null">{{ opt.label }}</button>
                  </li>
                </ul>
              </span>
              a message via
              <span class="create-campaign-modal__token-wrap">
                <button
                  type="button"
                  class="create-campaign-modal__token"
                  :class="{ 'is-placeholder': !aiChannel }"
                  @click.stop="toggleMenu('channel')"
                >
                  {{ channelTokenLabel }}
                  <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
                    <path d="M3 4.2L6 7.2 9 4.2" stroke="currentColor" stroke-width="1.5"
                      stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                </button>
                <ul v-if="openMenu === 'channel'" class="create-campaign-modal__menu">
                  <li v-for="opt in channelOptions" :key="opt.value">
                    <button type="button" @click="aiChannel = opt.value; openMenu = null">{{ opt.label }}</button>
                  </li>
                </ul>
              </span>
              when
              <span class="create-campaign-modal__token-wrap">
                <button
                  type="button"
                  class="create-campaign-modal__token"
                  :class="{ 'is-placeholder': !aiTrigger }"
                  @click.stop="toggleMenu('trigger')"
                >
                  {{ triggerTokenLabel }}
                  <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
                    <path d="M3 4.2L6 7.2 9 4.2" stroke="currentColor" stroke-width="1.5"
                      stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                </button>
                <ul v-if="openMenu === 'trigger'" class="create-campaign-modal__menu">
                  <li v-for="opt in triggerOptions" :key="opt.value">
                    <button type="button" @click="aiTrigger = opt.value; openMenu = null">{{ opt.label }}</button>
                  </li>
                </ul>
              </span>
            </p>
            <div class="create-campaign-modal__ai-foot">
              <button type="button" class="create-campaign-modal__icon-btn" aria-label="Reset" @click="resetAi">
                <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
                  <path d="M3 8a5 5 0 019.5-2.2M13 8a5 5 0 01-9.5 2.2" stroke="currentColor"
                    stroke-width="1.5" stroke-linecap="round" />
                  <path d="M12.5 2.5V5.8H9.2M3.5 13.5V10.2H6.8" stroke="currentColor"
                    stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
                </svg>
              </button>
              <button type="button" class="create-campaign-modal__submit" aria-label="Create automation" @click="submitAi">
                <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
                  <path d="M8 12V4M4.5 7.5L8 4l3.5 3.5" stroke="currentColor" stroke-width="1.7"
                    stroke-linecap="round" stroke-linejoin="round" />
                </svg>
              </button>
            </div>
          </div>

          <div v-else class="create-campaign-modal__ai-box">
            <label class="is-sr-only" for="ai-write">Describe your automation</label>
            <textarea
              id="ai-write"
              v-model="aiText"
              class="create-campaign-modal__write"
              rows="4"
              placeholder="Describe the automation you want, e.g. Send a welcome email when someone joins my list"
              @keydown.enter.ctrl.prevent="submitAi"
            />
            <div class="create-campaign-modal__ai-foot">
              <button type="button" class="create-campaign-modal__icon-btn" aria-label="Reset" @click="resetAi">
                <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
                  <path d="M3 8a5 5 0 019.5-2.2M13 8a5 5 0 01-9.5 2.2" stroke="currentColor"
                    stroke-width="1.5" stroke-linecap="round" />
                  <path d="M12.5 2.5V5.8H9.2M3.5 13.5V10.2H6.8" stroke="currentColor"
                    stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
                </svg>
              </button>
              <button type="button" class="create-campaign-modal__submit" aria-label="Create automation" @click="submitAi">
                <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
                  <path d="M8 12V4M4.5 7.5L8 4l3.5 3.5" stroke="currentColor" stroke-width="1.7"
                    stroke-linecap="round" stroke-linejoin="round" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </section>

      <section class="create-campaign-modal__section">
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
      aiMode: 'guided',
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
        { id: 'whatsapp', label: 'WhatsApp', badge: '' },
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
      this.aiOpen = true;
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
