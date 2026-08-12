<template>
  <section class="analytics content relative bv-page analytics-brevo">
    <marketing-subnav />

    <header class="analytics-brevo__header">
      <div>
        <h1 class="analytics-brevo__title">
          {{ $t('analytics.title') }}
        </h1>
        <p class="analytics-brevo__lead">Track opens, clicks, and performance across your campaigns.</p>
      </div>
    </header>
    <div v-if="serverConfig.privacy.disable_tracking"
      class="notification is-info">
      {{ $t('analytics.trackingDisabled') }}
    </div>
    <div v-else-if="showAnonymousNotice"
      class="notification is-info">
      {{ $t('analytics.nonIndividualTracking') }}
    </div>
    <div v-if="trackingNotes" class="notification is-light">
      <strong>{{ $t('analytics.trackingNotes') }}:</strong> {{ trackingNotes }}
    </div>

    <form @submit.prevent="onSubmit" class="bv-page__card mb-4">
      <div class="columns">
        <div class="column is-6">
          <b-field :label="$t('globals.terms.campaigns')" label-position="on-border">
            <b-taginput v-model="form.campaigns" :data="queriedCampaigns" name="campaigns" ellipsis icon="tag-outline"
              :placeholder="$t('globals.terms.campaigns')" autocomplete :allow-new="false" :open-on-focus="true"
              :before-adding="isCampaignSelected" @typing="queryCampaigns" @focus="queryCampaigns" field="name"
              :loading="isSearchLoading" />
          </b-field>
        </div>

        <div class="column is-5">
          <div class="columns">
            <div class="column is-6">
              <b-field data-cy="from" :label="$t('analytics.fromDate')" label-position="on-border">
                <b-datetimepicker v-model="form.from" icon="calendar-clock" :timepicker="{ hourFormat: '24' }"
                  :datetime-formatter="formatDateTime" @input="onFromDateChange" />
              </b-field>
            </div>
            <div class="column is-6">
              <b-field data-cy="to" :label="$t('analytics.toDate')" label-position="on-border">
                <b-datetimepicker v-model="form.to" icon="calendar-clock" :timepicker="{ hourFormat: '24' }"
                  :datetime-formatter="formatDateTime" @input="onToDateChange" />
              </b-field>
            </div>
          </div>
        </div>

        <div class="column is-1">
          <b-button native-type="submit" type="is-dark" icon-left="magnify" :disabled="form.campaigns.length === 0"
            data-cy="btn-search" />
        </div>
        <div class="column is-1" v-if="reportURL">
          <b-tooltip :label="$t('campaigns.downloadReport')" type="is-dark">
            <b-button tag="a" :href="reportURL" target="_blank" rel="noopener noreferrer" type="is-light"
              icon-left="cloud-download-outline" data-cy="btn-download-report" />
          </b-tooltip>
        </div>
      </div>
    </form>

    <b-tabs v-model="activeTab" class="mt-5" :animated="false" @input="onTabChange">
      <b-tab-item :label="$t('analytics.overview')" value="overview">
        <section class="charts">
          <div class="chart" v-for="(v, k) in charts" :key="k">
            <div class="columns">
              <div class="column is-9">
                <b-loading v-if="v.loading" :active="v.loading" :is-full-page="false" />
                <h4 v-if="v.data !== null">
                  {{ v.name }}
                  <span class="has-text-grey-light">({{ $utils.niceNumber(counts[k]) }})</span>
                </h4>
                <chart :type="v.type" v-if="!v.loading && v.data" :data="v.data" :on-click="v.onClick" />
              </div>
              <div class="column is-2 donut-container">
                <chart type="donut" v-if="!v.loading && v.donutData" :data="v.donutData" />
              </div>
            </div>
          </div>
        </section>
      </b-tab-item>

      <b-tab-item :label="$t('campaigns.views')" value="views" :disabled="!hasCampaigns">
        <analytics-records-table typ="views" :campaigns="form.campaigns" :from="form.from" :to="form.to"
          :anonymous="showAnonymousNotice" :key="recordsKey" />
      </b-tab-item>

      <b-tab-item :label="$t('campaigns.clicks')" value="clicks" :disabled="!hasCampaigns">
        <analytics-records-table typ="clicks" :campaigns="form.campaigns" :from="form.from" :to="form.to"
          :anonymous="showAnonymousNotice" :key="recordsKey" />
      </b-tab-item>

      <b-tab-item :label="$t('analytics.links')" value="links" :disabled="!hasCampaigns">
        <section class="charts">
          <div class="chart">
            <b-loading v-if="charts.links.loading" :active="charts.links.loading" :is-full-page="false" />
            <h4 v-if="charts.links.data">
              {{ charts.links.name }}
              <span class="has-text-grey-light">({{ $utils.niceNumber(counts.links) }})</span>
            </h4>
            <chart type="bar" v-if="!charts.links.loading && charts.links.data" :data="charts.links.data"
              :on-click="charts.links.onClick" />
          </div>
        </section>
        <b-table v-if="linkRows.length" :data="linkRows" striped hoverable class="mt-4">
          <b-table-column field="url" :label="$t('analytics.links')" v-slot="props">
            <a :href="props.row.url" target="_blank" rel="noopener noreferrer">{{ props.row.url }}</a>
          </b-table-column>
          <b-table-column field="count" :label="$t('analytics.count')" v-slot="props" numeric>
            {{ $utils.niceNumber(props.row.count) }}
          </b-table-column>
        </b-table>
      </b-tab-item>
    </b-tabs>
  </section>
</template>

<script>
import dayjs from 'dayjs';
import Vue from 'vue';
import { mapState } from 'vuex';
import { colors } from '../constants';
import Chart from '../components/Chart.vue';
import AnalyticsRecordsTable from '../components/AnalyticsRecordsTable.vue';
import MarketingSubnav from '../components/MarketingSubnav.vue';

const chartColorRed = '#ee7d5b';
const chartColors = [
  colors.primary,
  '#FFB50D',
  '#41AC9C',
  chartColorRed,
  '#7FC7BC',
  '#3a82d6',
  '#688ED9',
  '#FFC43D',
];

export default Vue.extend({
  components: {
    Chart,
    AnalyticsRecordsTable,
    MarketingSubnav,
  },

  data() {
    return {
      activeTab: 'overview',
      recordsKey: 0,
      isSearchLoading: false,
      queriedCampaigns: [],
      linkRows: [],

      counts: {
        views: 0,
        clicks: 0,
        bounces: 0,
        links: 0,
      },
      urls: [],
      charts: {
        views: {
          name: this.$t('campaigns.views'),
          type: 'line',
          data: null,
          fn: this.$api.getCampaignViewCounts,
          chartFn: this.makeCharts,
          loading: false,
        },
        clicks: {
          name: this.$t('campaigns.clicks'),
          type: 'line',
          data: null,
          fn: this.$api.getCampaignClickCounts,
          chartFn: this.makeCharts,
          loading: false,
        },
        bounces: {
          name: this.$t('globals.terms.bounces'),
          type: 'line',
          data: null,
          fn: this.$api.getCampaignBounceCounts,
          chartFn: this.makeCharts,
          donutColor: chartColorRed,
          loading: false,
        },
        links: {
          name: this.$t('analytics.links'),
          type: 'bar',
          data: null,
          loading: false,
          fn: this.$api.getCampaignLinkCounts,
          chartFn: this.makeLinksChart,
          onClick: this.onLinkClick,
        },
      },

      form: {
        campaigns: [],
        from: null,
        to: null,
      },
    };
  },

  computed: {
    ...mapState(['serverConfig']),

    hasCampaigns() {
      return this.form.campaigns.length > 0;
    },

    showAnonymousNotice() {
      if (this.serverConfig.privacy.disable_tracking) {
        return false;
      }
      if (!this.serverConfig.privacy.individual_tracking) {
        return true;
      }
      return this.form.campaigns.some((c) => {
        const tc = c.trackingConfig || c.tracking_config || {};
        return tc.individualTracking === false || tc.individual_tracking === false;
      });
    },

    trackingNotes() {
      const notes = this.form.campaigns
        .map((c) => {
          const tc = c.trackingConfig || c.tracking_config || {};
          return tc.notes ? `${c.name.replace(/^#\d+: /, '')}: ${tc.notes}` : '';
        })
        .filter(Boolean);
      return notes.join(' · ');
    },

    reportURL() {
      if (this.form.campaigns.length !== 1) {
        return '';
      }
      const c = this.form.campaigns[0];
      if (!['finished', 'paused', 'running'].includes(c.status) || !c.sent) {
        return '';
      }
      return `/api/campaigns/${c.id}/report`;
    },
  },

  methods: {
    onTabChange() {
      if (this.activeTab !== 'overview' && this.hasCampaigns) {
        this.recordsKey += 1;
      }
    },

    onFromDateChange() {
      if (this.form.from > this.form.to) {
        this.form.to = dayjs(this.form.from).add(7, 'day').toDate();
      }
    },

    onToDateChange() {
      if (this.form.from > this.form.to) {
        this.form.from = dayjs(this.form.to).add(-7, 'day').toDate();
      }
    },

    formatDateTime(s) {
      return dayjs(s).format('YYYY-MM-DD HH:mm');
    },

    isCampaignSelected(camp) {
      return !this.form.campaigns.find(({ id }) => id === camp.id);
    },

    makeLinksChart(typ, camps, data) {
      this.linkRows = data.map((l) => ({ url: l.url, count: l.count }));
      const labels = data.map((l) => {
        try {
          this.urls.push(l.url);
          const u = new URL(l.url);
          if (l.url.length > 80) {
            return `${u.hostname}${u.pathname.substr(0, 50)}..`;
          }
          return u.hostname + u.pathname;
        } catch {
          return l.url;
        }
      });

      return {
        points: {
          labels,
          datasets: [{ data: data.map((l) => l.count), backgroundColor: chartColors }],
        },
        donut: null,
      };
    },

    makeCharts(typ, campaigns, data) {
      const camps = campaigns.reduce((obj, c) => {
        const out = { ...obj };
        out[c.id] = c;
        return out;
      }, {});
      const campIDs = Object.keys(camps);
      const lines = campIDs.map((id, n) => {
        const cId = parseInt(id, 10);
        const points = data.filter((item) => item.campaignId === cId);
        return {
          label: camps[id].name,
          data: points.map((item) => ({ x: this.formatDateTime(item.timestamp), y: item.count })),
          borderColor: chartColors[n % chartColors.length],
          borderWidth: 2,
          pointHoverBorderWidth: 5,
          pointBorderWidth: 0.5,
        };
      });

      const labels = [];
      const points = campIDs.map((id) => {
        labels.push(camps[id].name);
        const cId = parseInt(id, 10);
        return data.reduce((a, item) => (item.campaignId === cId ? a + item.count : a), 0);
      });

      return {
        points: { datasets: lines },
        donut: { labels, datasets: [{ data: points, backgroundColor: chartColors, borderWidth: 6 }] },
      };
    },

    onSubmit() {
      this.$router.push({
        query: {
          id: this.form.campaigns.map((c) => c.id),
          from: dayjs(this.form.from).unix(),
          to: dayjs(this.form.to).unix(),
        },
      });
      this.fetchAll();
    },

    queryCampaigns(q) {
      this.isSearchLoading = true;
      this.$api.getCampaigns({
        query: q,
        order_by: 'created_at',
        order: 'DESC',
      }).then((data) => {
        this.isSearchLoading = false;
        this.queriedCampaigns = data.results.map((c) => {
          const camp = c;
          camp.name = `#${c.id}: ${c.name}`;
          return camp;
        });
      });
    },

    getData(typ, camps) {
      this.charts[typ].loading = true;
      this.charts[typ].fn({
        id: camps.map((c) => c.id),
        from: this.form.from,
        to: this.form.to,
      }).then((data) => {
        this.counts[typ] = data.reduce((sum, d) => sum + d.count, 0);
        const { points, donut } = this.charts[typ].chartFn(typ, camps, data);
        this.charts[typ].data = points;
        this.charts[typ].donutData = donut;
        this.charts[typ].loading = false;
      });
    },

    fetchAll() {
      if (!this.hasCampaigns) {
        return;
      }
      Object.keys(this.charts).forEach((k) => {
        this.charts[k].data = null;
        this.charts[k].donutData = null;
        this.getData(k, this.form.campaigns);
      });
      this.recordsKey += 1;
    },

    onLinkClick(e) {
      const bars = e.chart.getElementsAtEventForMode(e, 'nearest', { intersect: true }, true);
      if (bars.length > 0) {
        window.open(this.urls[bars[0].index], '_blank', 'noopener noreferrer');
      }
    },

    loadCampaignsFromRoute() {
      const ids = this.$utils.parseQueryIDs(this.$route.query.id);
      if (ids.length === 0) {
        return Promise.resolve();
      }
      this.isSearchLoading = true;
      return Promise.allSettled(ids.map((id) => this.$api.getCampaign(id))).then((data) => {
        data.forEach((d) => {
          if (d.status !== 'fulfilled') {
            return;
          }
          const camp = d.value;
          camp.name = `#${camp.id}: ${camp.name}`;
          if (!this.form.campaigns.find((c) => c.id === camp.id)) {
            this.form.campaigns.push(camp);
          }
        });
        this.isSearchLoading = false;
        this.fetchAll();
      });
    },
  },

  created() {
    const now = dayjs().set('hour', 23).set('minute', 59).set('seconds', 0);
    const weekAgo = now.subtract(7, 'day').set('hour', 0).set('minute', 0);
    const from = this.$route.query.from ? dayjs.unix(this.$route.query.from) : weekAgo;
    const to = this.$route.query.to ? dayjs.unix(this.$route.query.to) : now;
    this.form.from = from.toDate();
    this.form.to = to.toDate();
  },

  mounted() {
    this.loadCampaignsFromRoute();
  },
});
</script>
