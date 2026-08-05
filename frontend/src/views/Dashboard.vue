<template>
  <section class="home-brevo">
    <b-loading v-if="isLoading" :active="isLoading" :is-full-page="false" />

    <header class="home-brevo__hello">
      <h1 class="home-brevo__greeting">Hello {{ displayName }}</h1>
    </header>

    <section class="home-brevo__planner">
      <div class="home-brevo__calendar">
        <div class="home-brevo__cal-head">
          <button type="button" class="home-brevo__cal-nav" @click="shiftMonth(-1)" aria-label="Previous month">
            <b-icon icon="chevron-left" size="is-small" />
          </button>
          <span>{{ calendarLabel }}</span>
          <button type="button" class="home-brevo__cal-nav" @click="shiftMonth(1)" aria-label="Next month">
            <b-icon icon="chevron-right" size="is-small" />
          </button>
        </div>
        <div class="home-brevo__cal-weekdays">
          <span v-for="d in weekdays" :key="d">{{ d }}</span>
        </div>
        <div class="home-brevo__cal-grid">
          <button
            v-for="(cell, idx) in calendarCells"
            :key="idx"
            type="button"
            class="home-brevo__cal-day"
            :class="{
              'is-empty': !cell.day,
              'is-selected': cell.day && isSameDay(cell.date, selectedDate),
              'is-today': cell.day && isSameDay(cell.date, today),
              'has-activity': cell.day && activityDates.has(cell.key),
            }"
            :disabled="!cell.day"
            @click="selectDate(cell.date)"
          >
            <span v-if="cell.day">{{ cell.day }}</span>
          </button>
        </div>
      </div>

      <div class="home-brevo__planned">
        <div class="home-brevo__planned-head">
          <h2>Planned for {{ plannedLabel }}</h2>
          <b-dropdown position="is-bottom-left">
            <template #trigger>
              <button type="button" class="home-brevo__btn-dark">
                + Create
                <b-icon icon="arrow-down" size="is-small" />
              </button>
            </template>
            <b-dropdown-item
              v-if="$can('campaigns:manage')"
              @click.native="goCreateCampaign"
            >
              Campaign
            </b-dropdown-item>
            <b-dropdown-item @click.native="openTaskForm()">
              Task
            </b-dropdown-item>
          </b-dropdown>
        </div>

        <div v-if="plannedItems.length === 0" class="home-brevo__empty">
          Nothing planned on this day.
        </div>
        <div v-else class="home-brevo__campaign-list">
          <button
            v-for="task in plannedTasks"
            :key="`task-${task.id}`"
            type="button"
            class="home-brevo__campaign-row is-task"
            @click="openTaskForm(task)"
          >
            <div class="home-brevo__campaign-main">
              <div class="home-brevo__campaign-title">{{ task.name }}</div>
              <div class="home-brevo__campaign-meta">
                #{{ task.id }} · {{ task.time || '' }}
                <span class="home-brevo__status" :class="task.highPriority ? 'is-paused' : 'is-scheduled'">
                  <span class="dot" /> {{ task.highPriority ? 'High priority' : 'Task' }}
                </span>
                <span class="home-brevo__pill">{{ typeLabel(task.type) }}</span>
              </div>
            </div>
          </button>
          <router-link
            v-for="c in plannedCampaigns"
            :key="`planned-${c.id}`"
            :to="{ name: 'campaign', params: { id: c.id }, query: { view: 'report' } }"
            class="home-brevo__campaign-row"
          >
            <div class="home-brevo__campaign-main">
              <div class="home-brevo__campaign-title">{{ c.name }}</div>
              <div class="home-brevo__campaign-meta">
                #{{ c.id }} · {{ campaignSentLabel(c) }}
                <span class="home-brevo__status" :class="`is-${c.status}`">
                  <span class="dot" /> {{ statusLabel(c.status) }}
                </span>
                <span class="home-brevo__pill">Email</span>
              </div>
            </div>
            <div class="home-brevo__campaign-stats">
              <div>
                <span class="lbl">Open rate</span>
                <strong>{{ openRate(c) }}</strong>
              </div>
              <div>
                <span class="lbl">Click rate</span>
                <strong>{{ clickRate(c) }}</strong>
              </div>
              <div>
                <span class="lbl">Conversions</span>
                <strong>0</strong>
              </div>
            </div>
          </router-link>
        </div>
        <div class="home-brevo__plant" aria-hidden="true">🌿</div>
      </div>
    </section>

    <section class="home-brevo__grid">
      <article class="home-brevo__card">
        <header class="home-brevo__card-head">
          <h2>Your last campaigns</h2>
          <router-link
            v-if="$can('campaigns:manage')"
            :to="{ name: 'campaign', params: { id: 'new' } }"
            class="home-brevo__link"
          >
            Create a campaign
          </router-link>
        </header>
        <div v-if="recentCampaigns.length === 0" class="home-brevo__empty">
          No campaigns yet.
        </div>
        <div v-else class="home-brevo__campaign-list">
          <router-link
            v-for="c in recentCampaigns"
            :key="`recent-${c.id}`"
            :to="{ name: 'campaign', params: { id: c.id }, query: { view: 'report' } }"
            class="home-brevo__campaign-row is-compact"
          >
            <div class="home-brevo__campaign-main">
              <div class="home-brevo__campaign-title">{{ c.name }}</div>
              <div class="home-brevo__campaign-meta">
                #{{ c.id }} · {{ campaignSentLabel(c) }}
                <span class="home-brevo__status" :class="`is-${c.status}`">
                  <span class="dot" /> {{ statusLabel(c.status) }}
                </span>
                <span class="home-brevo__pill">Email</span>
              </div>
            </div>
            <div class="home-brevo__campaign-stats">
              <div>
                <span class="lbl">Open rate</span>
                <strong>{{ openRate(c) }}</strong>
              </div>
              <div>
                <span class="lbl">Click rate</span>
                <strong>{{ clickRate(c) }}</strong>
              </div>
              <div>
                <span class="lbl">Conversions</span>
                <strong>0</strong>
              </div>
            </div>
          </router-link>
        </div>
        <footer class="home-brevo__card-foot">
          <router-link :to="{ name: 'campaigns' }" class="home-brevo__link">
            Go to Campaigns
          </router-link>
        </footer>
      </article>

      <div class="home-brevo__stack">
        <article class="home-brevo__card home-brevo__contacts">
          <header class="home-brevo__card-head">
            <h2>Your contacts</h2>
            <router-link
              v-if="$can('subscribers:get_all', 'subscribers:get')"
              :to="{ name: 'subscribers' }"
              class="home-brevo__link"
            >
              Add contact
            </router-link>
          </header>
          <div class="home-brevo__contacts-body">
            <div>
              <div class="home-brevo__big-num">{{ $utils.niceNumber(totalContacts) }}</div>
              <div class="home-brevo__muted">Total contacts</div>
            </div>
            <div class="home-brevo__icon-badge" aria-hidden="true">
              <b-icon icon="account-plus-outline" />
            </div>
          </div>
          <footer class="home-brevo__card-foot">
            <router-link :to="{ name: 'subscribers' }" class="home-brevo__link">
              Go to Contacts
            </router-link>
          </footer>
        </article>

        <article class="home-brevo__card home-brevo__contacts is-secondary">
          <div class="home-brevo__contacts-body">
            <div>
              <div class="home-brevo__big-num">{{ $utils.niceNumber(newContacts) }}</div>
              <div class="home-brevo__muted">New contacts over the last 30 days</div>
            </div>
            <div class="home-brevo__icon-badge is-mint" aria-hidden="true">
              <b-icon icon="account-plus" />
            </div>
          </div>
        </article>
      </div>
    </section>

    <section class="home-brevo__cta">
      <div class="home-brevo__cta-art" aria-hidden="true">✉️</div>
      <div class="home-brevo__cta-body">
        <h2>Deliver the right message at the right time</h2>
        <p>
          Strengthen connections with timely and tailored messages, from welcome
          emails to follow-ups and order confirmation campaigns.
        </p>
        <router-link
          v-if="$can('campaigns:manage')"
          :to="{ name: 'campaign', params: { id: 'new' } }"
          class="home-brevo__btn-dark"
        >
          Create campaign
        </router-link>
      </div>
    </section>

    <b-modal :active.sync="isTaskFormVisible" :width="560" scroll="keep" class="task-form-modal">
      <task-form
        v-if="isTaskFormVisible"
        :data="editingTask"
        :default-date="selectedDate.toDate()"
        @close="isTaskFormVisible = false"
        @finished="onTaskFinished"
      />
    </b-modal>
  </section>
</template>

<script>
import dayjs from 'dayjs';
import Vue from 'vue';
import { mapState } from 'vuex';
import TaskForm, { loadTasks } from '../components/TaskForm.vue';

export default Vue.extend({
  name: 'Home',

  components: { TaskForm },

  data() {
    return {
      isLoading: true,
      campaigns: [],
      tasks: [],
      totalContacts: 0,
      newContacts: 0,
      today: dayjs().startOf('day'),
      selectedDate: dayjs().startOf('day'),
      viewMonth: dayjs().startOf('month'),
      weekdays: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
      isTaskFormVisible: false,
      editingTask: null,
    };
  },

  computed: {
    ...mapState(['profile']),

    displayName() {
      const p = this.profile || {};
      const name = (p.name || p.username || '').trim();
      if (!name) return 'there';
      return name.split(/\s+/)[0];
    },

    calendarLabel() {
      return this.viewMonth.format('MMMM YYYY');
    },

    plannedLabel() {
      return this.selectedDate.format('MMMM D, YYYY');
    },

    activityDates() {
      const set = new Set();
      this.campaigns.forEach((c) => {
        const d = this.campaignDay(c);
        if (d) set.add(d.format('YYYY-MM-DD'));
      });
      this.tasks.forEach((t) => {
        if (t.dueDate) set.add(t.dueDate);
      });
      return set;
    },

    calendarCells() {
      const start = this.viewMonth.startOf('month');
      let pad = start.day() - 1;
      if (pad < 0) pad = 6;
      const daysInMonth = this.viewMonth.daysInMonth();
      const cells = [];
      for (let i = 0; i < pad; i += 1) {
        cells.push({ day: null, date: null, key: `e-${i}` });
      }
      for (let d = 1; d <= daysInMonth; d += 1) {
        const date = this.viewMonth.date(d).startOf('day');
        cells.push({ day: d, date, key: date.format('YYYY-MM-DD') });
      }
      while (cells.length % 7 !== 0) {
        cells.push({ day: null, date: null, key: `t-${cells.length}` });
      }
      return cells;
    },

    plannedCampaigns() {
      return this.campaigns.filter((c) => {
        const d = this.campaignDay(c);
        return d && this.isSameDay(d, this.selectedDate);
      }).slice(0, 5);
    },

    plannedTasks() {
      const key = this.selectedDate.format('YYYY-MM-DD');
      return this.tasks.filter((t) => t.dueDate === key && t.status !== 'done');
    },

    plannedItems() {
      return [...this.plannedTasks, ...this.plannedCampaigns];
    },

    recentCampaigns() {
      return [...this.campaigns]
        .sort((a, b) => {
          const da = this.campaignDay(a);
          const db = this.campaignDay(b);
          if (!da && !db) return 0;
          if (!da) return 1;
          if (!db) return -1;
          return db.valueOf() - da.valueOf();
        })
        .slice(0, 5);
    },
  },

  methods: {
    typeLabel(type) {
      const map = {
        todo: 'To do', call: 'Call', email: 'Email', meeting: 'Meeting',
      };
      return map[type] || 'Task';
    },

    openTaskForm(task) {
      this.editingTask = task && task.id ? { ...task } : null;
      this.isTaskFormVisible = true;
    },

    goCreateCampaign() {
      this.$router.push({ name: 'campaign', params: { id: 'new' } });
    },

    onTaskFinished() {
      this.tasks = loadTasks();
    },

    fetchData() {
      this.isLoading = true;
      this.tasks = loadTasks();

      const countsP = this.$api.getDashboardCounts().then((data) => {
        this.totalContacts = (data.subscribers && data.subscribers.total) || 0;
      }).catch(() => {
        this.totalContacts = 0;
      });

      const campaignsP = this.$api.getCampaigns({
        page: 1,
        per_page: 50,
        order_by: 'created_at',
        order: 'desc',
        no_body: true,
        query: '',
      }).then((data) => {
        this.campaigns = (data && data.results) || [];
      }).catch(() => {
        this.campaigns = [];
      });

      const since = dayjs().subtract(30, 'day').format('YYYY-MM-DD');
      const newSubsP = this.$api.getSubscribers({
        page: 1,
        per_page: 1,
        query: `subscribers.created_at >= '${since}'`,
      }).then((data) => {
        this.newContacts = (data && data.total) || 0;
      }).catch(() => {
        this.newContacts = 0;
      });

      Promise.all([countsP, campaignsP, newSubsP]).finally(() => {
        this.isLoading = false;
      });
    },

    campaignDay(c) {
      const raw = c.startedAt || c.sendAt || c.createdAt;
      if (!raw) return null;
      return dayjs(raw).startOf('day');
    },

    isSameDay(a, b) {
      if (!a || !b) return false;
      return dayjs(a).isSame(dayjs(b), 'day');
    },

    selectDate(date) {
      if (!date) return;
      this.selectedDate = dayjs(date).startOf('day');
    },

    shiftMonth(delta) {
      this.viewMonth = this.viewMonth.add(delta, 'month').startOf('month');
    },

    statusLabel(status) {
      if (status === 'finished') return 'Sent';
      if (status === 'running') return 'Sending';
      if (status === 'scheduled') return 'Scheduled';
      if (status === 'paused') return 'Paused';
      if (status === 'cancelled') return 'Cancelled';
      if (status === 'draft') return 'Draft';
      return status || '—';
    },

    campaignSentLabel(c) {
      const raw = c.startedAt || c.sendAt || c.createdAt;
      if (!raw) return 'Not sent yet';
      const d = dayjs(raw);
      const prefix = ['finished', 'running'].includes(c.status) ? 'Sent on' : 'On';
      return `${prefix} ${d.format('MMM D, YYYY')} at ${d.format('h:mm A')}`;
    },

    metricVal(obj, keys) {
      for (let i = 0; i < keys.length; i += 1) {
        const v = obj && obj[keys[i]];
        if (v !== undefined && v !== null && v !== '') {
          const n = Number(v);
          if (!Number.isNaN(n)) {
            return n;
          }
        }
      }
      return 0;
    },

    rate(num, den) {
      if (!den || den <= 0) return '0.00%';
      const pct = (num / den) * 100;
      return `${pct.toFixed(2)}%`;
    },

    openRate(c) {
      const sent = this.metricVal(c, ['sent', 'toSend', 'to_send']);
      const opens = this.metricVal(c, ['views', 'opens', 'openCount']);
      return this.rate(opens, sent);
    },

    clickRate(c) {
      const sent = this.metricVal(c, ['sent', 'toSend', 'to_send']);
      const clicks = this.metricVal(c, ['clicks', 'clickCount']);
      return this.rate(clicks, sent);
    },
  },

  created() {
    this.$root.$on('page.refresh', this.fetchData);
  },

  destroyed() {
    this.$root.$off('page.refresh', this.fetchData);
  },

  mounted() {
    this.fetchData();
  },
});
</script>
