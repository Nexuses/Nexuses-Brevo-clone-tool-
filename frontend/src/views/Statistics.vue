<template>
  <section class="crm-page statistics-page statistics-brevo">
    <marketing-subnav />

    <header class="statistics-brevo__header">
      <div>
        <h1 class="statistics-brevo__title">Statistics</h1>
        <p class="statistics-brevo__lead">
          Overview of sends, opens, clicks, and bounce performance.
        </p>
      </div>
      <div class="stats-date-range">
        <b-datepicker
          v-model="dateRange"
          range
          placeholder="Select date range"
          icon="calendar-clock"
          :editable="false"
          @input="onDateChange"
        />
      </div>
    </header>

    <section class="stats-summary">
      <div class="stats-summary__head">
        <h2>Summary</h2>
        <span class="stats-summary__hint">
          Automated opens and clicks included
        </span>
      </div>

      <div class="stats-summary__grid" v-if="!loading">
        <div class="stats-metric">
          <div class="stats-metric__row">
            <span class="stats-metric__lbl">Sent</span>
            <strong class="stats-metric__val">{{ $utils.formatNumber(summary.sent) }}</strong>
          </div>
          <div class="stats-metric__row">
            <span class="stats-metric__lbl">Soft + Hard Bounces</span>
            <strong class="stats-metric__val is-sm">{{ $utils.formatNumber(summary.bounces) }}</strong>
          </div>
        </div>

        <div class="stats-metric">
          <div class="stats-metric__row">
            <span class="stats-metric__lbl">Recipients</span>
            <strong class="stats-metric__val">{{ $utils.formatNumber(summary.recipients) }}</strong>
          </div>
          <div class="stats-metric__row">
            <span class="stats-metric__lbl">Replies</span>
            <strong class="stats-metric__val is-sm">0</strong>
          </div>
        </div>

        <div class="stats-metric">
          <div class="stats-metric__row">
            <span class="stats-metric__lbl">Opens</span>
            <strong class="stats-metric__val">{{ $utils.formatNumber(summary.opens) }}</strong>
          </div>
          <div class="stats-metric__row">
            <span class="stats-metric__lbl">Open Rate</span>
            <strong class="stats-metric__val is-sm">{{ summary.openRate }}</strong>
          </div>
        </div>

        <div class="stats-metric">
          <div class="stats-metric__row">
            <span class="stats-metric__lbl">Clicks</span>
            <strong class="stats-metric__val">{{ $utils.formatNumber(summary.clicks) }}</strong>
          </div>
          <div class="stats-metric__row">
            <span class="stats-metric__lbl">Click Rate</span>
            <strong class="stats-metric__val is-sm">{{ summary.clickRate }}</strong>
          </div>
        </div>

        <div class="stats-metric">
          <div class="stats-metric__row">
            <span class="stats-metric__lbl">Unsubscribes</span>
            <strong class="stats-metric__val">0</strong>
          </div>
          <div class="stats-metric__row">
            <span class="stats-metric__lbl">Unsubscription Rate</span>
            <strong class="stats-metric__val is-sm">0%</strong>
          </div>
        </div>
      </div>
      <b-loading :active="loading" :is-full-page="false" />
    </section>

    <section class="stats-campaigns">
      <div class="stats-campaigns__head">
        <h2>Email Campaigns</h2>
        <a
          v-if="exportHref"
          class="stats-export"
          :href="exportHref"
          download="campaign-statistics.csv"
        >
          <b-icon icon="cloud-download-outline" size="is-small" />
          Export (.csv file)
        </a>
      </div>

      <div class="crm-card-table">
        <table class="table is-fullwidth is-hoverable">
          <thead>
            <tr>
              <th>ID</th>
              <th>Name</th>
              <th>Recipients</th>
              <th>Opened</th>
              <th>Clicked</th>
              <th>Unsubscribed</th>
              <th>Complained</th>
              <th>Bounces</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="c in filteredCampaigns" :key="c.id">
              <td class="has-text-grey">#{{ c.id }}</td>
              <td>
                <router-link
                  class="crm-link"
                  :to="{ name: 'campaign', params: { id: c.id }, query: { view: 'report' } }"
                >
                  {{ c.name }}
                </router-link>
              </td>
              <td>
                <strong>{{ $utils.formatNumber(recipients(c)) }}</strong>
                <span class="stats-pct">100%</span>
              </td>
              <td>
                <strong>{{ $utils.formatNumber(c.views || 0) }}</strong>
                <span class="stats-pct">{{ pct(c.views, c.sent) }}</span>
                <router-link
                  class="stats-details"
                  :to="{ name: 'campaign', params: { id: c.id }, query: { view: 'report' } }"
                >
                  Details
                </router-link>
              </td>
              <td class="is-clicks">
                <strong>{{ $utils.formatNumber(c.clicks || 0) }}</strong>
                <span class="stats-pct">{{ pct(c.clicks, c.sent) }}</span>
              </td>
              <td class="is-warn">
                <strong>0</strong>
                <span class="stats-pct">0%</span>
              </td>
              <td>
                <strong>0</strong>
                <span class="stats-pct">0%</span>
              </td>
              <td class="is-warn">
                <strong>{{ $utils.formatNumber(c.bounces || 0) }}</strong>
                <span class="stats-pct">{{ pct(c.bounces, c.sent) }}</span>
              </td>
            </tr>
            <tr v-if="!loading && filteredCampaigns.length === 0">
              <td colspan="8" class="has-text-centered has-text-grey py-5">
                No campaigns in this date range.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </section>
</template>

<script>
import dayjs from 'dayjs';
import MarketingSubnav from '../components/MarketingSubnav.vue';

export default {
  name: 'Statistics',
  components: { MarketingSubnav },

  data() {
    const to = new Date();
    const from = dayjs().subtract(30, 'day').toDate();
    return {
      loading: false,
      campaigns: [],
      dateRange: [from, to],
    };
  },

  computed: {
    filteredCampaigns() {
      if (!this.dateRange || this.dateRange.length < 2 || !this.dateRange[0] || !this.dateRange[1]) {
        return this.campaigns;
      }
      const from = dayjs(this.dateRange[0]).startOf('day');
      const to = dayjs(this.dateRange[1]).endOf('day');
      return this.campaigns.filter((c) => {
        const d = dayjs(c.updatedAt || c.createdAt);
        return d.isAfter(from) && d.isBefore(to);
      });
    },

    summary() {
      const list = this.filteredCampaigns;
      const sent = list.reduce((n, c) => n + (c.sent || 0), 0);
      const recipients = list.reduce((n, c) => n + (c.toSend || c.sent || 0), 0);
      const opens = list.reduce((n, c) => n + (c.views || 0), 0);
      const clicks = list.reduce((n, c) => n + (c.clicks || 0), 0);
      const bounces = list.reduce((n, c) => n + (c.bounces || 0), 0);
      return {
        sent,
        recipients,
        opens,
        clicks,
        bounces,
        openRate: this.pct(opens, sent),
        clickRate: this.pct(clicks, sent),
      };
    },

    exportHref() {
      if (this.filteredCampaigns.length === 0) return null;
      const header = ['ID', 'Name', 'Recipients', 'Opened', 'Clicked', 'Bounces'];
      const rows = this.filteredCampaigns.map((c) => [
        c.id,
        `"${(c.name || '').replace(/"/g, '""')}"`,
        this.recipients(c),
        c.views || 0,
        c.clicks || 0,
        c.bounces || 0,
      ].join(','));
      const csv = [header.join(','), ...rows].join('\n');
      return `data:text/csv;charset=utf-8,${encodeURIComponent(csv)}`;
    },
  },

  methods: {
    recipients(c) {
      return c.toSend || c.sent || 0;
    },

    pct(num, den) {
      if (!den || den <= 0) return '0%';
      const p = (num / den) * 100;
      return Number.isInteger(p) ? `${p}%` : `${p.toFixed(2)}%`;
    },

    onDateChange() {
      // filtering is reactive via computed
    },

    async loadCampaigns() {
      this.loading = true;
      try {
        const data = await this.$api.getCampaigns({
          page: 1,
          per_page: 'all',
          order_by: 'updated_at',
          order: 'desc',
          no_body: true,
        });
        this.campaigns = (data && data.results) || [];
      } catch (e) {
        this.campaigns = [];
      } finally {
        this.loading = false;
      }
    },
  },

  mounted() {
    this.loadCampaigns();
  },
};
</script>
