<template>
  <section class="lists lists-brevo">
    <crm-subnav />
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
      <div class="lists-brevo__folder-wrap" ref="folderWrap">
        <button
          type="button"
          class="lists-brevo__folder-chip"
          :class="{ 'is-open': folderFilterOpen }"
          aria-haspopup="listbox"
          :aria-expanded="folderFilterOpen ? 'true' : 'false'"
          @click.stop="toggleFolderFilter"
        >
          {{ folderChipLabel }}
          <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
            <path
              v-if="folderFilterOpen"
              d="M9 7.5L6 4.5L3 7.5"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
            <path
              v-else
              d="M3 4.5L6 7.5L9 4.5"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
        </button>

        <div
          v-show="folderFilterOpen"
          class="lists-brevo__folder-panel"
          @click.stop
        >
          <div class="lists-brevo__folder-panel-scroll">
            <button
              type="button"
              class="lists-brevo__folder-option"
              :class="{ 'is-active': !queryParams.folderId }"
              @click="setFolderFilter(null)"
            >
              <span class="lists-brevo__folder-option-name">All folders</span>
              <span class="lists-brevo__folder-option-meta">{{ listCountLabel(allListsCount) }}</span>
            </button>
            <button
              v-for="folder in folders"
              :key="folder.id"
              type="button"
              class="lists-brevo__folder-option"
              :class="{ 'is-active': queryParams.folderId === folder.id }"
              @click="setFolderFilter(folder.id)"
            >
              <span class="lists-brevo__folder-option-name">{{ folder.name }}</span>
              <span class="lists-brevo__folder-option-meta">{{ listCountLabel(folderCounts[folder.id] || 0) }}</span>
            </button>
          </div>
          <div class="lists-brevo__folder-menu-foot">
            <button type="button" class="lists-brevo__folder-create" @click="showCreateFolderDrawer">
              <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
                <circle cx="8" cy="8" r="6.25" stroke="currentColor" stroke-width="1.4" />
                <path d="M8 5.2v5.6M5.2 8h5.6" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" />
              </svg>
              Create a new folder
            </button>
          </div>
        </div>
      </div>

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

    <div class="bv-table-card" :class="{ 'has-selection': hasSelection }">
      <bv-select-bar
        :selected="numSelectedLists"
        :total="lists.total || 0"
        noun="lists"
        :all-selected="bulk.all"
        @clear="clearSelection"
        @select-all="selectAllLists"
      >
        <button
          v-if="$can('lists:manage_all')"
          type="button"
          class="bv-select-bar__action"
          data-cy="btn-delete-lists"
          @click="deleteLists"
        >
          Delete
        </button>
      </bv-select-bar>
    <b-table
      :data="lists.results"
      :loading="loading.listsFull"
      @check-all="onTableCheck"
      @check="onTableCheck"
      :checked-rows.sync="bulk.checked"
      hoverable
      checkable
      :mobile-cards="false"
      default-sort="createdAt"
      default-sort-direction="desc"
      :paginated="false"
      :current-page="queryParams.page"
      :per-page="lists.perPage"
      :total="lists.total"
      backend-sorting
      @sort="onSort"
      class="lists-brevo__table"
    >
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

      <b-table-column v-slot="props" field="type" label="Folder" header-class="cy-type">
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
        <b-dropdown
          position="is-bottom-left"
          class="bv-action-menu lists-brevo__actions-dd"
          :mobile-modal="false"
        >
          <button
            slot="trigger"
            type="button"
            class="lists-brevo__kebab"
            aria-label="Actions"
          >
            <span aria-hidden="true" /><span aria-hidden="true" /><span aria-hidden="true" />
          </button>

          <b-dropdown-item
            v-if="$can('lists:manage') || $canList(props.row.id, 'list:manage')"
            @click="showEditForm(props.row)"
          >
            <b-icon icon="cog-outline" size="is-small" />
            List settings
          </b-dropdown-item>
          <b-dropdown-item
            v-if="$can('lists:manage_all')"
            @click="duplicateList(props.row)"
          >
            <b-icon icon="content-copy" size="is-small" />
            Duplicate
          </b-dropdown-item>
          <b-dropdown-item
            v-if="$can('subscribers:import')"
            @click="importToList(props.row)"
          >
            <b-icon icon="cloud-upload-outline" size="is-small" />
            Import contacts to the list
          </b-dropdown-item>
          <b-dropdown-item
            v-if="$can('lists:manage') || $canList(props.row.id, 'list:manage')"
            @click="renameList(props.row)"
          >
            <b-icon icon="pencil-outline" size="is-small" />
            Rename
          </b-dropdown-item>
          <b-dropdown-item
            v-if="$can('lists:manage') || $canList(props.row.id, 'list:manage')"
            @click="openMoveDrawer(props.row)"
          >
            <b-icon icon="arrow-all" size="is-small" />
            Move
          </b-dropdown-item>
          <hr class="bv-action-menu__sep" />
          <b-dropdown-item
            v-if="$can('lists:manage') || $canList(props.row.id, 'list:manage')"
            class="bv-action-menu__danger"
            @click="deleteList(props.row)"
          >
            <b-icon icon="trash-can-outline" size="is-small" />
            Delete list
          </b-dropdown-item>
        </b-dropdown>
      </b-table-column>

      <template #empty v-if="!loading.listsFull">
        <empty-placeholder
          label="No lists yet"
          description="Create a list to organize contacts for campaigns and imports."
        />
      </template>
    </b-table>

    <div class="contacts-brevo__pager lists-brevo__pager">
      <div class="contacts-brevo__pager-meta">
        <span>{{ listsPagerSummary }}</span>
        <span class="contacts-brevo__pager-pages">Page {{ queryParams.page }} of {{ listsTotalPages }}</span>
      </div>
      <div class="contacts-brevo__pager-controls">
        <button
          type="button"
          class="contacts-brevo__pager-btn"
          :disabled="queryParams.page <= 1"
          @click="onPageChange(queryParams.page - 1)"
        >
          Previous
        </button>
        <button
          type="button"
          class="contacts-brevo__pager-btn"
          :disabled="queryParams.page >= listsTotalPages"
          @click="onPageChange(queryParams.page + 1)"
        >
          Next
        </button>
      </div>
    </div>
    </div>

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

          <div class="lists-drawer__field">
            <span class="lists-drawer__label">
              Folder <em aria-hidden="true">*</em>
            </span>
            <div class="lists-drawer__folder-wrap" ref="createFolderWrap">
              <button
                type="button"
                class="lists-drawer__folder-trigger"
                :class="{ 'is-open': createFolderDdOpen, 'is-placeholder': !createForm.folderId }"
                aria-haspopup="listbox"
                @click.stop="toggleCreateFolderDd"
              >
                {{ createFolderTriggerLabel }}
                <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
                  <path d="M3 4.5L6 7.5L9 4.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
                </svg>
              </button>
              <div
                v-show="createFolderDdOpen"
                class="lists-brevo__folder-panel lists-brevo__folder-panel--drawer"
                @click.stop
              >
                <div class="lists-brevo__folder-panel-scroll">
                  <button
                    v-for="folder in folders"
                    :key="folder.id"
                    type="button"
                    class="lists-brevo__folder-option"
                    :class="{ 'is-active': createForm.folderId === folder.id }"
                    @click="selectCreateFolder(folder.id)"
                  >
                    <span class="lists-brevo__folder-option-name">{{ folder.name }}</span>
                    <span class="lists-brevo__folder-option-meta">{{ listCountLabel(folderCounts[folder.id] || 0) }}</span>
                  </button>
                  <p v-if="folders.length === 0" class="lists-brevo__folder-empty">
                    No folders yet. Create one first.
                  </p>
                </div>
                <div class="lists-brevo__folder-menu-foot">
                  <button type="button" class="lists-brevo__folder-create" @click="showCreateFolderDrawer">
                    <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
                      <circle cx="8" cy="8" r="6.25" stroke="currentColor" stroke-width="1.4" />
                      <path d="M8 5.2v5.6M5.2 8h5.6" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" />
                    </svg>
                    Create a new folder
                  </button>
                </div>
              </div>
            </div>
          </div>

          <footer class="contacts-drawer__foot">
            <button type="button" class="contacts-drawer__cancel" @click="closeCreateDrawer">Cancel</button>
            <button
              type="submit"
              class="contacts-drawer__save lists-drawer__save"
              :disabled="isCreating || !createForm.name.trim() || !createForm.folderId"
            >
              Create list
            </button>
          </footer>
        </form>
      </aside>
    </div>

    <!-- Create a folder — Brevo side drawer -->
    <div
      class="contacts-drawer lists-drawer lists-folder-drawer"
      :class="{ 'is-open': isFolderDrawerOpen }"
      :aria-hidden="isFolderDrawerOpen ? 'false' : 'true'"
    >
      <button type="button" class="contacts-drawer__backdrop" aria-label="Close" @click="closeFolderDrawer" />
      <aside class="contacts-drawer__panel" role="dialog" aria-modal="true" aria-label="Create a folder">
        <header class="contacts-drawer__head">
          <h2>Create a folder</h2>
          <button type="button" class="contacts-drawer__close" aria-label="Close" @click="closeFolderDrawer">
            <svg width="14" height="14" viewBox="0 0 14 14" fill="none" aria-hidden="true">
              <path d="M3 3l8 8M11 3L3 11" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" />
            </svg>
          </button>
        </header>
        <form class="contacts-drawer__body lists-drawer__body" @submit.prevent="createFolderFromDrawer">
          <p class="lists-folder-drawer__hint">
            Organize your lists efficiently by creating folders. Choose a clear, specific name for easy retrieval.
          </p>
          <label class="lists-drawer__field">
            <span class="lists-drawer__label">
              Name of the folder <em aria-hidden="true">*</em>
            </span>
            <div class="lists-drawer__input-wrap">
              <input
                v-model="folderForm.name"
                type="text"
                placeholder="Enter a folder name"
                maxlength="50"
                required
                ref="folderName"
              />
              <span class="lists-drawer__counter">{{ folderForm.name.length }}/50</span>
            </div>
          </label>
          <footer class="contacts-drawer__foot">
            <button type="button" class="contacts-drawer__cancel" @click="closeFolderDrawer">Cancel</button>
            <button
              type="submit"
              class="contacts-drawer__save lists-drawer__save"
              :disabled="!folderForm.name.trim()"
            >
              Create folder
            </button>
          </footer>
        </form>
      </aside>
    </div>
    <!-- Move list — folder picker drawer -->
    <div
      class="contacts-drawer lists-drawer lists-folder-drawer"
      :class="{ 'is-open': isMoveDrawerOpen }"
      :aria-hidden="isMoveDrawerOpen ? 'false' : 'true'"
    >
      <button type="button" class="contacts-drawer__backdrop" aria-label="Close" @click="closeMoveDrawer" />
      <aside class="contacts-drawer__panel" role="dialog" aria-modal="true" aria-label="Move list">
        <header class="contacts-drawer__head">
          <h2>Move list</h2>
          <button type="button" class="contacts-drawer__close" aria-label="Close" @click="closeMoveDrawer">
            <svg width="14" height="14" viewBox="0 0 14 14" fill="none" aria-hidden="true">
              <path d="M3 3l8 8M11 3L3 11" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" />
            </svg>
          </button>
        </header>
        <div class="contacts-drawer__body lists-drawer__body">
          <p class="lists-folder-drawer__hint" v-if="movingList">
            Choose a folder for “{{ movingList.name }}”.
          </p>
          <div class="lists-brevo__folder-panel lists-brevo__folder-panel--static">
            <button
              v-for="folder in folders"
              :key="folder.id"
              type="button"
              class="lists-brevo__folder-option"
              :class="{ 'is-active': moveFolderId === folder.id }"
              @click="moveFolderId = folder.id"
            >
              <span class="lists-brevo__folder-option-name">{{ folder.name }}</span>
              <span class="lists-brevo__folder-option-meta">{{ listCountLabel(folderCounts[folder.id] || 0) }}</span>
            </button>
            <p v-if="folders.length === 0" class="lists-brevo__folder-empty">
              No folders yet. Create one first.
            </p>
          </div>
          <footer class="contacts-drawer__foot">
            <button type="button" class="contacts-drawer__cancel" @click="closeMoveDrawer">Cancel</button>
            <button
              type="button"
              class="contacts-drawer__save lists-drawer__save"
              :disabled="!moveFolderId || isMoving"
              @click="confirmMoveList"
            >
              Move
            </button>
          </footer>
        </div>
      </aside>
    </div>
  </section>
</template>

<script>
import Vue from 'vue';
import { mapState } from 'vuex';
import dayjs from 'dayjs';
import EmptyPlaceholder from '../components/EmptyPlaceholder.vue';
import CrmSubnav from '../components/CrmSubnav.vue';
import BvSelectBar from '../components/BvSelectBar.vue';
import ListForm from './ListForm.vue';

const FOLDERS_KEY = 'nexuses.crm.listFolders';
const FOLDER_TAG_PREFIX = 'folder-';

function folderTag(id) {
  return `${FOLDER_TAG_PREFIX}${id}`;
}

function parseFolderIdFromTags(tags) {
  if (!Array.isArray(tags)) return null;
  const match = tags.find((t) => typeof t === 'string' && t.startsWith(FOLDER_TAG_PREFIX));
  if (!match) return null;
  const id = Number(match.slice(FOLDER_TAG_PREFIX.length));
  return Number.isFinite(id) ? id : null;
}

export default Vue.extend({
  components: {
    ListForm,
    EmptyPlaceholder,
    CrmSubnav,
    BvSelectBar,
  },

  data() {
    return {
      curItem: null,
      isEditing: false,
      isFormVisible: false,
      isCreateDrawerOpen: false,
      isFolderDrawerOpen: false,
      isMoveDrawerOpen: false,
      isCreating: false,
      isMoving: false,
      folderFilterOpen: false,
      createFolderDdOpen: false,
      movingList: null,
      moveFolderId: null,
      createForm: {
        name: '',
        folderId: null,
      },
      folderForm: {
        name: '',
      },
      folders: [],
      allListsForCounts: [],
      lists: { results: [], total: 0, perPage: 20 },
      queryParams: {
        page: 1,
        query: '',
        folderId: null,
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
    loadFolders() {
      try {
        const raw = localStorage.getItem(FOLDERS_KEY);
        const parsed = raw ? JSON.parse(raw) : [];
        this.folders = Array.isArray(parsed) ? parsed : [];
      } catch (e) {
        this.folders = [];
      }
    },

    saveFolders() {
      localStorage.setItem(FOLDERS_KEY, JSON.stringify(this.folders));
    },

    listCountLabel(n) {
      return `${n} ${n === 1 ? 'list' : 'lists'}`;
    },

    folderById(id) {
      return this.folders.find((f) => f.id === id) || null;
    },

    folderLabel(row) {
      const id = parseFolderIdFromTags(row.tags);
      if (!id) return '—';
      const folder = this.folderById(id);
      return folder ? folder.name : '—';
    },

    setFolderFilter(folderId) {
      this.queryParams.folderId = folderId;
      this.queryParams.page = 1;
      this.folderFilterOpen = false;
      this.getLists();
    },

    toggleFolderFilter() {
      this.createFolderDdOpen = false;
      this.folderFilterOpen = !this.folderFilterOpen;
    },

    toggleCreateFolderDd() {
      this.folderFilterOpen = false;
      this.createFolderDdOpen = !this.createFolderDdOpen;
    },

    selectCreateFolder(folderId) {
      this.createForm.folderId = folderId;
      this.createFolderDdOpen = false;
    },

    onDocClick(e) {
      const { folderWrap, createFolderWrap } = this.$refs;
      if (this.folderFilterOpen && folderWrap && !folderWrap.contains(e.target)) {
        this.folderFilterOpen = false;
      }
      if (this.createFolderDdOpen && createFolderWrap && !createFolderWrap.contains(e.target)) {
        this.createFolderDdOpen = false;
      }
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
      this.closeFolderDrawer();
      this.curItem = list;
      this.isEditing = true;
      this.isFormVisible = true;
    },

    showNewForm() {
      this.isEditing = false;
      this.isFormVisible = false;
      this.closeFolderDrawer();
      this.createForm = {
        name: '',
        folderId: this.folders.length === 1 ? this.folders[0].id : null,
      };
      this.isCreateDrawerOpen = true;
      this.$nextTick(() => {
        if (this.$refs.createName) this.$refs.createName.focus();
      });
    },

    closeCreateDrawer() {
      this.isCreateDrawerOpen = false;
    },

    showCreateFolderDrawer() {
      this.folderFilterOpen = false;
      this.createFolderDdOpen = false;
      this.folderForm = { name: '' };
      this.isFolderDrawerOpen = true;
      this.$nextTick(() => {
        if (this.$refs.folderName) this.$refs.folderName.focus();
      });
    },

    closeFolderDrawer() {
      this.isFolderDrawerOpen = false;
    },

    importToList(list) {
      this.$router.push({ name: 'import', query: { list_id: list.id } });
    },

    renameList(list) {
      this.$utils.prompt('Rename list', {
        placeholder: 'List name',
        value: list.name,
      }, (name) => {
        const trimmed = (name || '').trim();
        if (!trimmed || trimmed === list.name) return;
        this.$api.updateList({
          id: list.id,
          name: trimmed,
          type: list.type,
          optin: list.optin,
          status: list.status,
          tags: list.tags || [],
          description: list.description || '',
        }).then(() => {
          this.getLists();
          this.$utils.toast(this.$t('globals.messages.updated', { name: trimmed }));
        });
      });
    },

    duplicateList(list) {
      const folderId = parseFolderIdFromTags(list.tags);
      const tags = folderId ? [folderTag(folderId)] : (list.tags || []);
      this.$api.createList({
        name: `Copy of ${list.name}`,
        type: list.type || 'private',
        optin: list.optin || 'single',
        status: 'active',
        tags,
        description: list.description || '',
      }).then((data) => {
        this.refreshFolderCounts();
        this.getLists();
        this.$utils.toast(this.$t('globals.messages.created', { name: data.name }));
      });
    },

    openMoveDrawer(list) {
      this.movingList = list;
      this.moveFolderId = parseFolderIdFromTags(list.tags);
      this.isMoveDrawerOpen = true;
    },

    closeMoveDrawer() {
      this.isMoveDrawerOpen = false;
      this.movingList = null;
      this.moveFolderId = null;
    },

    confirmMoveList() {
      if (!this.movingList || !this.moveFolderId) return;
      const list = this.movingList;
      const otherTags = (list.tags || []).filter((t) => !String(t).startsWith(FOLDER_TAG_PREFIX));
      this.isMoving = true;
      this.$api.updateList({
        id: list.id,
        name: list.name,
        type: list.type,
        optin: list.optin,
        status: list.status,
        tags: [...otherTags, folderTag(this.moveFolderId)],
        description: list.description || '',
      }).then(() => {
        this.closeMoveDrawer();
        this.refreshFolderCounts();
        this.getLists();
        this.$utils.toast('List moved');
      }).finally(() => {
        this.isMoving = false;
      });
    },

    createFolderFromDrawer() {
      const name = (this.folderForm.name || '').trim();
      if (!name) return;

      const exists = this.folders.some((f) => f.name.toLowerCase() === name.toLowerCase());
      if (exists) {
        this.$utils.toast('A folder with this name already exists');
        return;
      }

      const nextId = this.folders.reduce((max, f) => Math.max(max, f.id), 0) + 1;
      this.folders.push({ id: nextId, name });
      this.saveFolders();
      this.createForm.folderId = nextId;
      this.closeFolderDrawer();
      this.$utils.toast(`Folder "${name}" created`);
    },

    createListFromDrawer() {
      const name = (this.createForm.name || '').trim();
      const { folderId } = this.createForm;
      if (!name || !folderId) return;

      this.isCreating = true;
      this.$api.createList({
        name,
        type: 'private',
        optin: 'single',
        status: 'active',
        tags: [folderTag(folderId)],
      }).then((data) => {
        this.closeCreateDrawer();
        this.refreshFolderCounts();
        this.getLists();
        this.$utils.toast(this.$t('globals.messages.created', { name: data.name }));
      }).finally(() => {
        this.isCreating = false;
      });
    },

    formFinished() {
      this.refreshFolderCounts();
      this.getLists();
    },

    onFormClose() {
      this.isFormVisible = false;
      if (this.$route.params.id) {
        this.$router.push({ name: 'lists' });
      }
    },

    refreshFolderCounts() {
      return this.$api.queryLists({
        page: 1,
        per_page: 'all',
        status: this.queryParams.status || 'active',
      }).then((resp) => {
        this.allListsForCounts = resp.results || [];
      }).catch(() => {
        this.allListsForCounts = [];
      });
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
      if (this.queryParams.folderId) {
        params.tag = folderTag(this.queryParams.folderId);
      }

      this.$api.queryLists(params).then((resp) => {
        this.lists = resp;
      });

      this.$api.getLists({ minimal: true, per_page: 'all', status: 'active' });
      this.refreshFolderCounts();
    },

    deleteList(list) {
      this.$utils.confirm(
        this.$t('lists.confirmDelete'),
        () => {
          this.$api.deleteList(list.id).then(() => {
            this.refreshFolderCounts();
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

    selectAllLists() {
      this.bulk.all = true;
      this.bulk.checked = [...(this.lists.results || [])];
    },

    clearSelection() {
      this.bulk.all = false;
      this.bulk.checked = [];
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
            this.refreshFolderCounts();
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

    hasSelection() {
      return this.bulk.checked.length > 0 || this.bulk.all;
    },

    allListsCount() {
      return this.allListsForCounts.length;
    },

    folderCounts() {
      const counts = {};
      this.folders.forEach((f) => {
        counts[f.id] = 0;
      });
      this.allListsForCounts.forEach((list) => {
        const id = parseFolderIdFromTags(list.tags);
        if (id && Object.prototype.hasOwnProperty.call(counts, id)) {
          counts[id] += 1;
        }
      });
      return counts;
    },

    listsTotalPages() {
      const per = this.lists.perPage || 20;
      const total = this.lists.total || 0;
      return Math.max(1, Math.ceil(total / per));
    },

    listsPagerSummary() {
      const total = this.lists.total || 0;
      const per = this.lists.perPage || 20;
      if (!total) return '0 lists';
      const start = ((this.queryParams.page - 1) * per) + 1;
      const end = Math.min(this.queryParams.page * per, total);
      return `Showing ${start}–${end} of ${total} lists`;
    },

    folderChipLabel() {
      if (this.queryParams.folderId) {
        const folder = this.folderById(this.queryParams.folderId);
        const count = this.folderCounts[this.queryParams.folderId] || 0;
        const label = folder ? folder.name : 'Folder';
        return `${label} (${this.listCountLabel(count)})`;
      }
      return `All folders (${this.listCountLabel(this.lists.total || this.allListsCount)})`;
    },

    createFolderTriggerLabel() {
      if (!this.createForm.folderId) return 'Select a folder';
      const folder = this.folderById(this.createForm.folderId);
      return folder ? folder.name : 'Select a folder';
    },
  },

  created() {
    this.loadFolders();
    this.$root.$on('page.refresh', this.getLists);
  },

  destroyed() {
    clearTimeout(this.searchDebounce);
    this.$root.$off('page.refresh', this.getLists);
    document.removeEventListener('click', this.onDocClick);
  },

  mounted() {
    document.addEventListener('click', this.onDocClick);
    this.getLists();
  },
});
</script>
