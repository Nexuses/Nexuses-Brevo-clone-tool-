<template>
  <section class="logs logs-brevo content relative bv-page">
    <header class="logs-brevo__header">
      <h1 class="logs-brevo__title">
        {{ $t('logs.title') }}
      </h1>
      <p class="logs-brevo__lead">Live application logs from the server.</p>
    </header>
    <div class="bv-page__card">
      <log-view :loading="loading.logs" :lines="lines" />
    </div>
  </section>
</template>

<script>
import Vue from 'vue';
import { mapState } from 'vuex';
import LogView from '../components/LogView.vue';

export default Vue.extend({
  components: {
    LogView,
  },

  data() {
    return {
      lines: [],
      pollId: null,
    };
  },

  methods: {
    getLogs() {
      this.$api.getLogs().then((data) => {
        this.lines = data;
      });
    },
  },

  computed: {
    ...mapState(['logs', 'loading']),
  },

  mounted() {
    this.getLogs();

    // Update the logs every 10 seconds.
    this.pollId = setInterval(() => this.getLogs(), 10000);
  },

  destroyed() {
    clearInterval(this.pollId);
  },
});
</script>
