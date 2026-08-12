<template>
  <section class="lists lists-brevo">
    <header class="lists-brevo__header">
      <div class="lists-brevo__title-row">
        <h1 class="lists-brevo__title">Lists</h1>
        <button
          v-if="$can('lists:manage_all')"
          type="button"
          class="lists-brevo__create"
          data-cy="btn-new"
          @click="showNewForm"
        >
          <span class="lists-brevo__create-plus" aria-hidden="true">+</span>
          Create a list
        </button>
      </div>
      <p class="lists-brevo__lead">
        This is where you organize your lists. Create, modify, and manage custom lists for targeted
        interactions, and keep them in folders for easy navigation.
      </p>
      <div class="lists-brevo__meta-row">
        <div class="lists-brevo__links">
          <a
            href="https://listmonk.app/docs/lists/"
            target="_blank"
            rel="noopener noreferrer"
            class="lists-brevo__link lists-brevo__link--doc"
          >
            Get started with Lists and Folders
            <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
              <path
                d="M5 2.5H2.5A1 1 0 001.5 3.5v6A1 1 0 002.5 10.5h6a1 1 0 001-1V7M7 1.5h3.5V5M5.5 6.5L10.5 1.5"
                stroke="currentColor"
                stroke-width="1.3"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </a>
          <a
            href="https://listmonk.app/docs/querying/#lists"
            target="_blank"
            rel="noopener noreferrer"
            class="lists-brevo__link lists-brevo__link--doc"
          >
            Lists vs Segments
            <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
              <path
                d="M5 2.5H2.5A1 1 0 001.5 3.5v6A1 1 0 002.5 10.5h6a1 1 0 001-1V7M7 1.5h3.5V5M5.5 6.5L10.5 1.5"
                stroke="currentColor"
                stroke-width="1.3"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </a>
        </div>
        <button
          v-if="queryParams.status !== 'archived'"
          type="button"
          class="lists-brevo__link lists-brevo__link--recalc"
          @click="getLists"
        >
          Recalculate now
        </button>
        <router-link
          v-else
          :to="{ name: 'lists' }"
          class="lists-brevo__link lists-brevo__link--recalc"
        >
          View all lists
        </router-link>
      </div>
    </header>

    <div class="lists-brevo__toolbar">
      <button type="button" class="lists-brevo__folder-chip">
        All folders ({{ lists.total || 0 }} lists)
        <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
          <path d="M3 4.5L6 7.5L9 4.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
      </button>
      <form class="lists-brevo__search" @submit.prevent="getLists">
        <b-field>
          <b-input
            v-model="queryParams.query"
            name="query"
            expanded
            icon="magnify"
            placeholder="Search a list name or ID"
            ref="query"
            data-cy="query"
            @keyup.native.enter="getLists"
            @input="onSearchInput"
          />
        </b-field>
      </form>
    </div>

    <b-table
      :data="lists.results"
      :loading="loading.listsFull"
      @check-all="onTableCheck"
      @check="onTableCheck"
      :checked-rows.sync="bulk.checked"
      hoverable
      :mobile-cards="false"
      default-sort="createdAt"
      default-sort-direction="desc"
      paginated
      backend-pagination
      pagination-position="bottom"
      @page-change="onPageChange"
      :current-page="queryParams.page"
      :per-page="lists.perPage"
      :total="lists.total"
      backend-sorting
      @sort="onSort"
      class="lists-brevo__table"
    >
      <template #top-left>
        <div class="actions" v-if="bulk.checked.length > 0">
          <a class="a" href="#" @click.prevent="deleteLists" data-cy="btn-delete-lists">
            <b-icon icon="trash-can-outline" size="is-small" /> {{ $t('globals.buttons.delete') }}
          </a>
          <span class="a">
            {{ $tc('globals.messages.numSelected', numSelectedLists, { num: numSelectedLists }) }}
          </span>
        </div>
      </template>

      <b-table-column
        v-slot="props"
        field="name"
        label="Lists"
        header-class="cy-name"
        sortable
        :td-attrs="$utils.tdID"
      >
        <a href="#" class="lists-brevo__name" @click.prevent.stop="openListContacts(props.row.id)">
          {{ props.row.name }}
        </a>
      </b-table-column>

      <b-table-column v-slot="props" field="id" label="ID" header-class="cy-id" sortable>
        <span class="lists-brevo__id">#{{ props.row.id }}</span>
      </b-table-column>

      <b-table-column v-slot="props" field="type" label="Folder" header-class="cy-type" sortable>
        <span class="lists-brevo__folder">{{ folderLabel(props.row) }}</span>
      </b-table-column>

      <b-table-column
        v-slot="props"
        field="subscriber_count"
        label="Contacts"
        header-class="cy-subscribers"
        numeric
        sortable
      >
        <a
          v-if="$can('subscribers:get_all', 'subscribers:get')"
          href="#"
          class="lists-brevo__count"
          @click.prevent.stop="openListContacts(props.row.id)"
        >
          {{ $utils.formatNumber(props.row.subscriberCount) }}
        </a>
        <span v-else class="lists-brevo__count">
          {{ $utils.formatNumber(props.row.subscriberCount) }}
        </span>
      </b-table-column>

      <b-table-column
        v-slot="props"
        field="created_at"
        label="Creation date"
        header-class="cy-created_at"
        sortable
      >
        <span class="lists-brevo__date">{{ formatCreated(props.row.createdAt) }}</span>
      </b-table-column>

      <b-table-column v-slot="props" label="Actions" cell-class="actions" align="right" width="70">
        <b-dropdown position="is-bottom-left" class="campaign-actions-menu">
          <template #trigger>
            <button type="button" class="lists-brevo__kebab" aria-label="Actions">
              <span aria-hidden="true" /><span aria-hidden="true" /><span aria-hidden="true" />
            </button>
          </template>
          <div class="campaign-actions-panel">
            <a
              v-if="$can('subscribers:get_all', 'subscribers:get')"
              href="#"
              class="campaign-action"
              aria-label="View contacts"
              @click.prevent="openListContacts(props.row.id)"
            >
              <b-tooltip label="View contacts" type="is-dark" position="is-left">
                <b-icon icon="account-multiple" />
              </b-tooltip>
            </a>
            <a
              v-if="$can('lists:manage') || $canList(props.row.id, 'list:manage')"
              href="#"
              class="campaign-action"
              aria-label="Edit"
              @click.prevent="showEditForm(props.row)"
            >
              <b-tooltip :label="$t('globals.buttons.edit')" type="is-dark" position="is-left">
                <b-icon icon="pencil-outline" />
              </b-tooltip>
            </a>
            <router-link
              v-if="$can('campaigns:manage')"
              :to="`/campaigns/new?list_id=${props.row.id}`"
              class="campaign-action"
              aria-label="Send campaign"
            >
              <b-tooltip :label="$t('lists.sendCampaign')" type="is-dark" position="is-left">
                <b-icon icon="rocket-launch-outline" />
              </b-tooltip>
            </router-link>
            <a
              v-if="$can('lists:manage') || $canList(props.row.id, 'list:manage')"
              href="#"
              class="campaign-action"
              aria-label="Delete"
              @click.prevent="deleteList(props.row)"
            >
              <b-tooltip :label="$t('globals.buttons.delete')" type="is-dark" position="is-left">
                <b-icon icon="trash-can-outline" />
              </b-tooltip>
            </a>
          </div>
        </b-dropdown>
      </b-table-column>

      <template #empty v-if="!loading.listsFull">
        <empty-placeholder />
      </template>
    </b-table>

    <!-- Edit list keeps modal -->
    <b-modal
      v-if="isEditing"
      scroll="keep"
      :aria-modal="true"
      :active.sync="isFormVisible"
      :width="600"
      @close="onFormClose"
    >
      <list-form :data="curItem" :is-editing="true" @finished="formFinished" />
    </b-modal>

    <!-- Create a list — Brevo side drawer -->
    <div
      class="contacts-drawer lists-drawer"
      :class="{ 'is-open': isCreateDrawerOpen }"
      :aria-hidden="isCreateDrawerOpen ? 'false' : 'true'"
    >
      <button type="button" class="contacts-drawer__backdrop" aria-label="Close" @click="closeCreateDrawer" />
      <aside class="contacts-drawer__panel" role="dialog" aria-modal="true" aria-label="Create a list">
        <header class="contacts-drawer__head">
          <h2>Create a list</h2>
          <button type="button" class="contacts-drawer__close" aria-label="Close" @click="closeCreateDrawer">
            <svg width="14" height="14" viewBox="0 0 14 14" fill="none" aria-hidden="true">
              <path d="M3 3l8 8M11 3L3 11" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" />
            </svg>
          </button>
        </header>
        <form class="contacts-drawer__body lists-drawer__body" @submit.prevent="createListFromDrawer">
          <label class="lists-drawer__field">
            <span class="lists-drawer__label">
              Name of the list <em aria-hidden="true">*</em>
            </span>
            <div class="lists-drawer__input-wrap">
              <input
                v-model="createForm.name"
                type="text"
                placeholder="Enter a list name"
                maxlength="255"
                required
                ref="createName"
              />
              <span class="lists-drawer__counter">{{ createForm.name.length }}/255</span>
            </div>
          </label>

          <label class="lists-drawer__field">
            <span class="lists-drawer__label">
              Folder <em aria-hidden="true">*</em>
            </span>
            <div class="lists-drawer__select-wrap">
              <select v-model="createForm.type" required>
                <option value="" disabled>Select a folder</option>
                <option value="private">Private</option>
                <option value="public">Public</option>
              </select>
            </div>
          </label>

          <footer class="contacts-drawer__foot">
            <button type="button" class="contacts-drawer__cancel" @click="closeCreateDrawer">Cancel</button>
            <button
              type="submit"
              class="contacts-drawer__save lists-drawer__save"
              :disabled="isCreating || !createForm.name.trim() || !createForm.type"
            >
              Create list
            </button>
          </footer>
        </form>
      </aside>
    </div>
  </section>
</template>

<script>
import Vue from 'vue';
import { mapState } from 'vuex';
import dayjs from 'dayjs';
import EmptyPlaceholder from '../components/EmptyPlaceholder.vue';
import ListForm from './ListForm.vue';

export default Vue.extend({
  components: {
    ListForm,
    EmptyPlaceholder,
  },

  data() {
    return {
      curItem: null,
      isEditing: false,
      isFormVisible: false,
      isCreateDrawerOpen: false,
      isCreating: false,
      createForm: {
        name: '',
        type: '',
      },
      lists: { results: [], total: 0, perPage: 20 },
      queryParams: {
        page: 1,
        query: '',
        orderBy: 'created_at',
        order: 'desc',
        status: this.$route.query.status || 'active',
      },
      bulk: {
        checked: [],
        all: false,
      },
      searchDebounce: null,
    };
  },

  methods: {
    folderLabel(row) {
      if (row.type === 'public') return 'Public';
      if (row.type === 'private') return 'Private';
      return row.type || '—';
    },

    formatCreated(raw) {
      if (!raw) return '—';
      const d = dayjs(raw);
      if (!d.isValid()) return '—';
      return d.format('MMM D, YYYY HH:mm');
    },

    onPageChange(p) {
      this.queryParams.page = p;
      this.getLists();
    },

    onSort(field, direction) {
      this.queryParams.orderBy = field;
      this.queryParams.order = direction;
      this.getLists();
    },

    onSearchInput() {
      clearTimeout(this.searchDebounce);
      if (!this.queryParams.query || this.queryParams.query.trim() === '') {
        this.queryParams.page = 1;
        this.getLists();
        return;
      }
      this.searchDebounce = setTimeout(() => {
        this.queryParams.page = 1;
        this.getLists();
      }, 250);
    },

    openListContacts(listID) {
      this.$router.push({ name: 'subscribers_list', params: { listID } });
    },

    showEditForm(list) {
      this.closeCreateDrawer();
      this.curItem = list;
      this.isEditing = true;
      this.isFormVisible = true;
    },

    showNewForm() {
      this.isEditing = false;
      this.isFormVisible = false;
      this.createForm = { name: '', type: '' };
      this.isCreateDrawerOpen = true;
      this.$nextTick(() => {
        if (this.$refs.createName) this.$refs.createName.focus();
      });
    },

    closeCreateDrawer() {
      this.isCreateDrawerOpen = false;
    },

    createListFromDrawer() {
      const name = (this.createForm.name || '').trim();
      const { type } = this.createForm;
      if (!name || !type) return;

      this.isCreating = true;
      this.$api.createList({
        name,
        type,
        optin: 'single',
        status: 'active',
        tags: [],
      }).then((data) => {
        this.closeCreateDrawer();
        this.getLists();
        this.$utils.toast(this.$t('globals.messages.created', { name: data.name }));
      }).finally(() => {
        this.isCreating = false;
      });
    },

    formFinished() {
      this.getLists();
    },

    onFormClose() {
      this.isFormVisible = false;
      if (this.$route.params.id) {
        this.$router.push({ name: 'lists' });
      }
    },

    getLists() {
      const cleanedQuery = (this.queryParams.query || '').replace(/[^\p{L}\p{N}\s]/gu, ' ').trim();
      const params = {
        page: this.queryParams.page,
        order_by: this.queryParams.orderBy,
        order: this.queryParams.order,
        status: this.queryParams.status,
      };
      if (cleanedQuery) {
        params.query = cleanedQuery;
      }

      this.$api.queryLists(params).then((resp) => {
        this.lists = resp;
      });

      this.$api.getLists({ minimal: true, per_page: 'all', status: 'active' });
    },

    deleteList(list) {
      this.$utils.confirm(
        this.$t('lists.confirmDelete'),
        () => {
          this.$api.deleteList(list.id).then(() => {
            this.getLists();
            this.$utils.toast(this.$t('globals.messages.deleted', { name: list.name }));
          });
        },
      );
    },

    onTableCheck() {
      if (this.bulk.checked.length !== this.lists.total) {
        this.bulk.all = false;
      }
    },

    deleteLists() {
      const name = this.$tc('globals.terms.list', this.numSelectedLists);

      const fn = () => {
        const params = {};
        if (!this.bulk.all && this.bulk.checked.length > 0) {
          params.id = this.bulk.checked.map((l) => l.id);
        } else {
          params.query = this.queryParams.query.replace(/[^\p{L}\p{N}\s]/gu, ' ');
          params.all = this.bulk.all;
        }

        this.$api.deleteLists(params)
          .then(() => {
            this.getLists();
            this.$utils.toast(this.$tc(
              'globals.messages.deletedCount',
              this.numSelectedLists,
              { num: this.numSelectedLists, name },
            ));
          });
      };

      this.$utils.confirm(this.$tc(
        'globals.messages.confirmDelete',
        this.numSelectedLists,
        { num: this.numSelectedLists, name: name.toLowerCase() },
      ), fn);
    },
  },

  computed: {
    ...mapState(['loading', 'settings']),

    numSelectedLists() {
      return this.bulk.all ? this.lists.total : this.bulk.checked.length;
    },
  },

  created() {
    this.$root.$on('page.refresh', this.getLists);
  },

  destroyed() {
    clearTimeout(this.searchDebounce);
    this.$root.$off('page.refresh', this.getLists);
  },

  mounted() {
    this.getLists();
  },
});
</script>
