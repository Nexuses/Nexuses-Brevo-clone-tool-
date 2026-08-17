<template>
  <section class="sender-add bv-page">
    <h1 class="sender-add__title">{{ isEdit ? 'Edit sender' : 'Add Sender' }}</h1>
    <p class="sender-add__lead">
      Specify what your recipients will see when they receive emails from this sender.
    </p>

    <form class="sender-add__form" @submit.prevent="onSave">
      <div class="sender-add__grid">
        <div class="sender-add__fields">
          <label class="sender-add__field">
            <span class="sender-add__label">From Name <em>*</em></span>
            <p class="sender-add__hint">
              Specify what your recipients will see when they receive emails from this sender.
            </p>
            <input
              v-model="form.name"
              class="sender-add__input"
              type="text"
              name="from_name"
              maxlength="200"
              placeholder="John Doe"
              data-cy="from-name"
              required
            >
          </label>

          <label class="sender-add__field">
            <span class="sender-add__label">From Email <em>*</em></span>
            <p class="sender-add__hint">
              From email is the sender email address from which your recipients will receive
              your emails. Format must be name@domain.com.
            </p>
            <input
              v-model="form.email"
              class="sender-add__input"
              type="email"
              name="from_email"
              maxlength="200"
              placeholder="john.doe@email.com"
              data-cy="from-email"
              required
            >
          </label>
        </div>

        <aside class="sender-add__preview" aria-hidden="true">
          <div class="phone-mock">
            <div class="phone-mock__bezel">
              <div class="phone-mock__notch" />
              <div class="phone-mock__screen">
                <div class="phone-mock__status">
                  <span>9:41</span>
                  <span class="phone-mock__dots">●●●</span>
                </div>
                <div class="phone-mock__inbox-title">Inbox</div>
                <div class="phone-mock__mail">
                  <div class="phone-mock__avatar">{{ initials }}</div>
                  <div class="phone-mock__meta">
                    <strong>{{ previewName }}</strong>
                    <span>{{ previewEmail }}</span>
                    <em>Your campaign subject line</em>
                  </div>
                </div>
                <div class="phone-mock__mail is-dim">
                  <div class="phone-mock__avatar is-grey">A</div>
                  <div class="phone-mock__meta">
                    <strong>Another sender</strong>
                    <span>hello@example.com</span>
                    <em>Earlier message</em>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </aside>
      </div>

      <footer class="sdd-wizard__footer">
        <button type="button" class="sdd-wizard__cancel" @click="onCancel">Cancel</button>
        <button
          type="submit"
          class="sdd-wizard__primary"
          data-cy="btn-add-sender"
          :disabled="loading.settings || !canSave"
        >
          {{ isEdit ? 'Save' : 'Add sender' }}
        </button>
      </footer>
    </form>
  </section>
</template>

<script>
import Vue from 'vue';
import { mapState } from 'vuex';

export default Vue.extend({
  data() {
    return {
      form: {
        name: '',
        email: '',
      },
    };
  },

  computed: {
    ...mapState(['loading']),

    isEdit() {
      return Boolean(this.$route.query.email);
    },

    canSave() {
      return Boolean((this.form.name || '').trim() && (this.form.email || '').trim());
    },

    previewName() {
      return (this.form.name || '').trim() || 'John Doe';
    },

    previewEmail() {
      return (this.form.email || '').trim() || 'john.doe@email.com';
    },

    initials() {
      const name = this.previewName;
      const parts = name.split(/\s+/).filter(Boolean);
      if (parts.length >= 2) {
        return (parts[0][0] + parts[1][0]).toUpperCase();
      }
      return name.slice(0, 1).toUpperCase();
    },
  },

  methods: {
    formatFrom(name, email) {
      const n = (name || '').trim();
      const e = (email || '').trim();
      if (!e) return '';
      return n ? `${n} <${e}>` : e;
    },

    onCancel() {
      this.$router.push({ name: 'trackingDomains', query: { tab: 'senders' } });
    },

    onSave() {
      if (!this.canSave) return;
      const value = this.formatFrom(this.form.name, this.form.email);
      this.$api.updateAppFromEmail(value).then(() => {
        this.$utils.toast(this.$t('globals.messages.updated', { name: 'Sender' }));
        this.$router.push({ name: 'trackingDomains', query: { tab: 'senders' } });
      });
    },
  },

  mounted() {
    this.form.name = this.$route.query.name || '';
    this.form.email = this.$route.query.email || '';
  },
});
</script>
