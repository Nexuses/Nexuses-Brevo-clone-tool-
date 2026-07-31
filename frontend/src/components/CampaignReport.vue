<template>
  <section class="campaign-report">
    <header class="campaign-report__top">
      <router-link :to="{ name: 'campaigns' }" class="campaign-report__back" aria-label="Back">
        <b-icon icon="arrow-left" />
      </router-link>

      <div class="campaign-report__thumb" aria-hidden="true">
        <div class="campaign-report__thumb-inner">
          <b-icon icon="email-outline" size="is-large" />
        </div>
      </div>

      <div class="campaign-report__heading">
        <h1 class="campaign-report__title">{{ campaign.name }}</h1>
        <p class="campaign-report__sub">
          #{{ campaign.id }}
          <span v-if="sentLabel"> · {{ sentLabel }}</span>
        </p>
        <div class="campaign-report__meta">
          <div>
            <span class="lbl">Subject</span>
            <strong>{{ campaign.subject || '—' }}</strong>
          </div>
          <div>
            <span class="lbl">From</span>
            <strong>{{ campaign.fromEmail || '—' }}</strong>
          </div>
          <div>
            <span class="lbl">Reply to</span>
            <strong>{{ campaign.replyTo || campaign.fromEmail || '—' }}</strong>
          </div>
        </div>
      </div>

      <div class="campaign-report__actions">
        <b-button
          v-if="canExport"
          tag="a"
          :href="`/api/campaigns/${campaign.id}/report`"
          target="_blank"
          rel="noopener noreferrer"
          type="is-light"
          icon-left="cloud-download-outline"
        >
          Export report
        </b-button>
        <b-button type="is-light" icon-left="pencil-outline" @click="$emit('edit')">
          Edit
        </b-button>
      </div>
    </header>

    <nav class="campaign-report__tabs">
      <button
        v-for="t in tabs"
        :key="t.id"
        type="button"
        class="campaign-report__tab"
        :class="{ 'is-active': activeTab === t.id }"
        @click="activeTab = t.id"
      >
        {{ t.label }}
      </button>
    </nav>

    <!-- Overview -->
    <div v-if="activeTab === 'overview'" class="campaign-report__panel">
      <div class="campaign-report__section-head">
        <h2>Campaign performance</h2>
        <span class="hint">Automated opens and clicks included.</span>
      </div>

      <div class="campaign-report__metrics">
        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Delivered</span>
            <router-link
              v-if="$can('subscribers:get_all', 'subscribers:get')"
              :to="{ name: 'subscribers' }"
              class="metric__view"
            >
              <b-icon icon="account-multiple-outline" size="is-small" /> View
            </router-link>
          </div>
          <div class="metric__value">{{ $utils.niceNumber(delivered) }}</div>
          <div class="metric__rate">Delivery rate <strong>{{ deliveryRate }}</strong></div>
        </div>

        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Opens</span>
            <a href="#" class="metric__view" @click.prevent="activeTab = 'opens'">
              <b-icon icon="account-multiple-outline" size="is-small" /> View
            </a>
          </div>
          <div class="metric__value">{{ $utils.niceNumber(opens) }}</div>
          <div class="metric__rate">Open rate <strong>{{ openRate }}</strong></div>
        </div>

        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Clicks</span>
            <a href="#" class="metric__view" @click.prevent="activeTab = 'clicks'">
              <b-icon icon="account-multiple-outline" size="is-small" /> View
            </a>
          </div>
          <div class="metric__value">{{ $utils.niceNumber(clicks) }}</div>
          <div class="metric__rate">Click-through rate <strong>{{ clickRate }}</strong></div>
        </div>

        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Conversions</span>
          </div>
          <div class="metric__value">0</div>
          <div class="metric__rate">Conversion rate <strong>0%</strong></div>
        </div>

        <div class="metric">
          <div class="metric__top">
            <span class="metric__label">Unsubscribes</span>
            <router-link
              :to="{ name: 'bounces', query: { campaign_id: campaign.id } }"
              class="metric__view"
            >
              <b-icon icon="account-multiple-outline" size="is-small" /> View
            </router-link>
          </div>
          <div class="metric__value">{{ $utils.niceNumber(unsubscribes) }}</div>
          <div class="metric__rate">Unsubscribe rate <strong>{{ unsubRate }}</strong></div>
        </div>
      </div>

      <h2 class="campaign-report__h2">Campaign audience</h2>
      <div class="campaign-report__audience">
        <div class="audience__head">
          <b-icon icon="plus-circle-outline" size="is-small" />
          Included lists
        </div>
        <div v-if="audienceLists.length === 0" class="audience__empty">
          No lists attached to this campaign.
        </div>
        <div
          v-for="l in audienceLists"
          :key="l.id"
          class="audience__row"
        >
          <div class="audience__name">#{{ l.id }} {{ l.name }}</div>
          <div class="audience__count">{{ $utils.niceNumber(l.subscriberCount || 0) }} contacts</div>
          <router-link
            :to="{ name: 'subscribers_list', params: { listID: l.id } }"
            class="audience__view"
          >
            <b-icon icon="account-multiple-outline" size="is-small" /> View
          </router-link>
        </div>
      </div>

      <h2 class="campaign-report__h2">Timeline</h2>
      <div class="campaign-report__timeline">
        <div v-for="(ev, i) in timeline" :key="i" class="tl-item">
          <div class="tl-icon">
            <b-icon :icon="ev.icon" size="is-small" />
          </div>
          <div class="tl-body">
            <div class="tl-title">{{ ev.title }}</div>
            <div class="tl-desc">{{ ev.description }}</div>
            <div class="tl-time">{{ ev.time }}</div>
          </div>
        </div>
        <div v-if="timeline.length === 0" class="audience__empty">
          No timeline events yet.
        </div>
      </div>
    </div>

    <!-- Other tabs -->
    <div v-else-if="activeTab === 'deliverability'" class="campaign-report__panel">
      <div class="campaign-report__simple-card">
        <div class="columns is-mobile">
          <div class="column">
            <div class="metric__label">Delivered</div>
            <div class="metric__value">{{ $utils.niceNumber(delivered) }}</div>
            <div class="metric__rate">{{ deliveryRate }}</div>
          </div>
          <div class="column">
            <div class="metric__label">Bounces</div>
            <div class="metric__value">{{ $utils.niceNumber(bounces) }}</div>
            <div class="metric__rate">
              <router-link :to="{ name: 'bounces', query: { campaign_id: campaign.id } }" class="metric__view">
                View bounces
              </router-link>
            </div>
          </div>
          <div class="column">
            <div class="metric__label">To send</div>
            <div class="metric__value">{{ $utils.niceNumber(toSend) }}</div>
          </div>
        </div>
      </div>
    </div>

    <div v-else-if="activeTab === 'opens'" class="campaign-report__panel">
      <div class="campaign-report__simple-card">
        <div class="metric__value">{{ $utils.niceNumber(opens) }}</div>
        <div class="metric__rate">Open rate {{ openRate }}</div>
        <p class="mt-4">
          <router-link
            :to="{ name: 'campaignAnalytics', query: { id: campaign.id } }"
            class="metric__view"
          >
            Open detailed analytics →
          </router-link>
        </p>
      </div>
    </div>

    <div v-else-if="activeTab === 'clicks'" class="campaign-report__panel">
      <div class="campaign-report__simple-card">
        <div class="metric__value">{{ $utils.niceNumber(clicks) }}</div>
        <div class="metric__rate">Click-through rate {{ clickRate }}</div>
        <p class="mt-4">
          <router-link
            :to="{ name: 'campaignAnalytics', query: { id: campaign.id } }"
            class="metric__view"
          >
            Open detailed analytics →
          </router-link>
        </p>
      </div>
    </div>

    <div v-else-if="activeTab === 'conversions'" class="campaign-report__panel">
      <div class="campaign-report__simple-card">
        <div class="metric__value">0</div>
        <div class="metric__rate">Conversion rate 0%</div>
        <p class="has-text-grey mt-3">Conversions are not tracked in this installation.</p>
      </div>
    </div>

    <div v-else-if="activeTab === 'unsubscribes'" class="campaign-report__panel">
      <div class="campaign-report__simple-card">
        <div class="metric__value">{{ $utils.niceNumber(unsubscribes) }}</div>
        <div class="metric__rate">Unsubscribe rate {{ unsubRate }}</div>
      </div>
    </div>
  </section>
</template>

<script>
import dayjs from 'dayjs';
import Vue from 'vue';
import { mapState } from 'vuex';

export default Vue.extend({
  name: 'CampaignReport',

  props: {
    campaign: { type: Object, required: true },
  },

  data() {
    return {
      activeTab: 'overview',
      tabs: [
        { id: 'overview', label: 'Overview' },
        { id: 'deliverability', label: 'Deliverability' },
        { id: 'opens', label: 'Opens' },
        { id: 'clicks', label: 'Clicks' },
        { id: 'conversions', label: 'Conversions' },
        { id: 'unsubscribes', label: 'Unsubscribes' },
      ],
    };
  },

  computed: {
    ...mapState(['lists']),

    canExport() {
      return ['finished', 'paused', 'running'].includes(this.campaign.status)
        && (this.campaign.sent || 0) > 0;
    },

    delivered() {
      return this.campaign.sent || 0;
    },

    toSend() {
      return this.campaign.toSend || this.campaign.sent || 0;
    },

    opens() {
      return this.campaign.views || 0;
    },

    clicks() {
      return this.campaign.clicks || 0;
    },

    bounces() {
      return this.campaign.bounces || 0;
    },

    unsubscribes() {
      // listmonk does not expose per-campaign unsubscribes on the campaign object.
      return 0;
    },

    deliveryRate() {
      return this.pct(this.delivered, this.toSend);
    },

    openRate() {
      return this.pct(this.opens, this.delivered);
    },

    clickRate() {
      return this.pct(this.clicks, this.delivered);
    },

    unsubRate() {
      return this.pct(this.unsubscribes, this.delivered);
    },

    sentLabel() {
      const c = this.campaign;
      const raw = c.startedAt || c.sendAt || c.updatedAt;
      if (!raw) return '';
      const d = dayjs(raw);
      if (c.status === 'finished' || c.status === 'running') {
        return `Sent on ${d.format('MMM D, YYYY h:mm A')}`;
      }
      if (c.status === 'scheduled' && c.sendAt) {
        return `Scheduled for ${dayjs(c.sendAt).format('MMM D, YYYY h:mm A')}`;
      }
      return `${this.$t(`campaigns.status.${c.status}`)} · ${d.format('MMM D, YYYY h:mm A')}`;
    },

    audienceLists() {
      const attached = this.campaign.lists || [];
      const all = (this.lists && this.lists.results) || [];
      return attached.map((l) => {
        const full = all.find((x) => x.id === l.id);
        return {
          id: l.id,
          name: l.name || (full && full.name) || `List ${l.id}`,
          subscriberCount: (full && full.subscriberCount) || l.subscriberCount || 0,
        };
      });
    },

    timeline() {
      const c = this.campaign;
      const events = [];
      const name = c.name || 'campaign';

      if (c.status === 'finished' || (c.status === 'running' && c.sent > 0)) {
        const when = c.updatedAt || c.startedAt;
        events.push({
          icon: 'send',
          title: 'Sending completed',
          description: `The campaign [${c.id}] ${name} has been sent.`,
          time: when ? dayjs(when).format('MMM D, YYYY h:mm A') : '',
        });
      }

      if (c.sendAt) {
        events.push({
          icon: 'calendar',
          title: 'Scheduled',
          description: `The campaign [${c.id}] ${name} has been scheduled for ${dayjs(c.sendAt).format('MMM D, YYYY h:mm A')}.`,
          time: c.createdAt ? dayjs(c.createdAt).format('MMM D, YYYY h:mm A') : '',
        });
      } else if (c.startedAt) {
        events.push({
          icon: 'rocket-launch-outline',
          title: 'Sending started',
          description: `The campaign [${c.id}] ${name} started sending.`,
          time: dayjs(c.startedAt).format('MMM D, YYYY h:mm A'),
        });
      }

      if (c.createdAt) {
        events.push({
          icon: 'pencil-outline',
          title: 'Created',
          description: `The campaign [${c.id}] ${name} was created.`,
          time: dayjs(c.createdAt).format('MMM D, YYYY h:mm A'),
        });
      }

      return events;
    },
  },

  methods: {
    pct(num, den) {
      if (!den || den <= 0) return '—';
      return `${((num / den) * 100).toFixed(2)}%`;
    },
  },

  mounted() {
    // Ensure lists are available for subscriber counts.
    if (!this.lists || !this.lists.results || this.lists.results.length === 0) {
      this.$api.getLists({ per_page: 'all', minimal: true }).catch(() => {});
    }
  },
});
</script>
