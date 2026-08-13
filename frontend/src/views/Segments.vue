<template>
  <section class="crm-page segments-page segments-brevo">
    <crm-subnav />

    <header class="segments-brevo__header">
      <div class="segments-brevo__header-main">
        <h1 class="segments-brevo__title">Segments</h1>
        <p class="segments-brevo__lead">
          This is where you organize your segments. Create, modify, and manage segments for targeted
          interactions, and keep them in folders for easy navigation.
        </p>
      </div>
      <button
        type="button"
        class="segments-brevo__create"
        data-cy="btn-new-segment"
        @click="createSegment"
      >
        <span class="segments-brevo__create-plus" aria-hidden="true">+</span>
        Create a segment
      </button>
    </header>

    <div class="crm-toolbar">
      <div class="crm-toolbar__folder">
        All folders ({{ filteredSegments.length }} segments)
      </div>
      <b-field>
        <b-input
          v-model="query"
          expanded
          placeholder="Search a segment name or ID"
          icon="magnify"
          data-cy="segment-search"
        />
      </b-field>
    </div>

    <div class="crm-card-table">
      <table class="table is-fullwidth is-hoverable">
        <thead>
          <tr>
            <th>Segment</th>
            <th>ID</th>
            <th>Folder</th>
            <th>Contacts</th>
            <th>Last edit</th>
            <th class="has-text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in filteredSegments" :key="s.id">
            <td>
              <a href="#" class="crm-link" @click.prevent="editSegment(s)">{{ s.name }}</a>
            </td>
            <td class="has-text-grey">#{{ s.id }}</td>
            <td>{{ s.folder }}</td>
            <td>{{ $utils.formatNumber(s.contacts) }}</td>
            <td>{{ formatDate(s.updatedAt) }}</td>
            <td class="has-text-right">
              <b-dropdown
                class="bv-action-menu"
                position="is-bottom-left"
                :mobile-modal="false"
              >
                <template #trigger>
                  <button type="button" class="campaign-actions-trigger" aria-label="Actions">
                    <span class="campaign-kebab" aria-hidden="true">
                      <span /><span /><span />
                    </span>
                  </button>
                </template>
                <b-dropdown-item @click="editSegment(s)">
                  <b-icon icon="pencil-outline" size="is-small" />
                  Edit
                </b-dropdown-item>
                <hr class="bv-action-menu__sep" />
                <b-dropdown-item class="bv-action-menu__danger" @click="deleteSegment(s)">
                  <b-icon icon="trash-can-outline" size="is-small" />
                  Delete
                </b-dropdown-item>
              </b-dropdown>
            </td>
          </tr>
          <tr v-if="filteredSegments.length === 0">
            <td colspan="6" class="has-text-centered has-text-grey py-5">
              No segments yet. Create one to get started.
            </td>
          </tr>
        </tbody>
      </table>
      <div class="crm-card-table__foot">
        <span>{{ filteredSegments.length }} of {{ segments.length }}</span>
      </div>
    </div>
  </section>
</template>

<script>
import dayjs from 'dayjs';
import CrmSubnav from '../components/CrmSubnav.vue';

const STORAGE_KEY = 'nexuses.crm.segments';

const DEFAULT_SEGMENTS = [
  {
    id: 1,
    name: 'Engaged contacts',
    folder: 'My segments',
    contacts: 0,
    updatedAt: new Date().toISOString(),
  },
  {
    id: 2,
    name: 'Recent signups',
    folder: 'My segments',
    contacts: 0,
    updatedAt: new Date().toISOString(),
  },
];

export default {
  name: 'Segments',
  components: { CrmSubnav },

  data() {
    return {
      query: '',
      segments: [],
    };
  },

  computed: {
    filteredSegments() {
      const q = this.query.trim().toLowerCase();
      if (!q) return this.segments;
      return this.segments.filter((s) => (
        s.name.toLowerCase().includes(q) || String(s.id).includes(q)
      ));
    },
  },

  methods: {
    formatDate(iso) {
      return dayjs(iso).format('DD/MM/YYYY');
    },

    load() {
      try {
        const raw = localStorage.getItem(STORAGE_KEY);
        this.segments = raw ? JSON.parse(raw) : [...DEFAULT_SEGMENTS];
      } catch (e) {
        this.segments = [...DEFAULT_SEGMENTS];
      }
    },

    save() {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(this.segments));
    },

    createSegment() {
      this.$utils.prompt('Create a segment', {
        placeholder: 'Segment name',
        value: '',
      }, (name) => {
        const trimmed = (name || '').trim();
        if (!trimmed) return;
        const nextId = this.segments.reduce((max, s) => Math.max(max, s.id), 0) + 1;
        this.segments.unshift({
          id: nextId,
          name: trimmed,
          folder: 'My segments',
          contacts: 0,
          updatedAt: new Date().toISOString(),
        });
        this.save();
        this.$utils.toast('Segment created');
      });
    },

    editSegment(s) {
      this.$utils.prompt('Edit segment', {
        placeholder: 'Segment name',
        value: s.name,
      }, (name) => {
        const trimmed = (name || '').trim();
        if (!trimmed) return;
        this.segments = this.segments.map((item) => (
          item.id === s.id
            ? { ...item, name: trimmed, updatedAt: new Date().toISOString() }
            : item
        ));
        this.save();
        this.$utils.toast('Segment updated');
      });
    },

    deleteSegment(s) {
      this.$utils.confirm(`Delete segment "${s.name}"?`, () => {
        this.segments = this.segments.filter((x) => x.id !== s.id);
        this.save();
        this.$utils.toast('Segment deleted');
      });
    },
  },

  mounted() {
    this.load();
  },
};
</script>
