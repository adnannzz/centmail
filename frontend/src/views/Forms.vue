<template>
  <section class="forms content relative">
    <header class="columns page-header">
      <div class="column is-10">
        <h1 class="title is-4">
          {{ $t('forms.title') }}
          <span v-if="forms.length > 0">({{ forms.length }})</span>
        </h1>
      </div>
      <div class="column has-text-right">
        <b-field v-if="$can('forms:manage')" expanded>
          <b-button expanded type="is-primary" icon-left="plus" class="btn-new" @click="showNewForm"
            data-cy="btn-new">
            {{ $t('globals.buttons.new') }}
          </b-button>
        </b-field>
      </div>
    </header>
    <p class="has-text-grey">
      {{ $t('forms.pageHelp') }}
    </p>
    <hr />

    <b-table :data="forms" :hoverable="true" :loading="loading.forms" default-sort="createdAt">
      <b-table-column v-slot="props" field="name" :label="$t('globals.fields.name')" :td-attrs="$utils.tdID" sortable>
        <a href="#" @click.prevent="showEditForm(props.row)">
          {{ props.row.name }}
        </a>
      </b-table-column>

      <b-table-column v-slot="props" field="lists" :label="$t('globals.terms.lists')">
        <span class="has-text-grey">{{ (props.row.listIds || []).length }}</span>
      </b-table-column>

      <b-table-column v-slot="props" field="createdAt" :label="$t('globals.fields.createdAt')" sortable>
        {{ $utils.niceDate(props.row.createdAt) }}
      </b-table-column>

      <b-table-column v-slot="props" cell-class="actions" align="right">
        <div>
          <a href="#" @click.prevent="showEmbed(props.row)" data-cy="btn-code" :aria-label="$t('forms.getCode')">
            <b-tooltip :label="$t('forms.getCode')" type="is-dark">
              <b-icon icon="code" size="is-small" />
            </b-tooltip>
          </a>
          <a href="#" @click.prevent="showEditForm(props.row)" data-cy="btn-edit"
            :aria-label="$t('globals.buttons.edit')">
            <b-tooltip :label="$t('globals.buttons.edit')" type="is-dark">
              <b-icon icon="pencil-outline" size="is-small" />
            </b-tooltip>
          </a>
          <a href="#" @click.prevent="deleteForm(props.row)" data-cy="btn-delete"
            :aria-label="$t('globals.buttons.delete')">
            <b-tooltip :label="$t('globals.buttons.delete')" type="is-dark">
              <b-icon icon="trash-can-outline" size="is-small" />
            </b-tooltip>
          </a>
        </div>
      </b-table-column>

      <template #empty v-if="!loading.forms">
        <empty-placeholder />
      </template>
    </b-table>

    <!-- Add / edit form modal -->
    <b-modal scroll="keep" :aria-modal="true" :active.sync="isFormVisible" :width="700">
      <form-edit :data="curItem" :is-editing="isEditing" @finished="formFinished" />
    </b-modal>

    <!-- Embed code modal -->
    <b-modal scroll="keep" :aria-modal="true" :active.sync="isEmbedVisible" :width="900">
      <form-embed :data="curItem" />
    </b-modal>
  </section>
</template>

<script>
import Vue from 'vue';
import { mapState } from 'vuex';
import EmptyPlaceholder from '../components/EmptyPlaceholder.vue';
import FormEdit from './FormEdit.vue';
import FormEmbed from './FormEmbed.vue';

export default Vue.extend({
  name: 'Forms',

  components: {
    EmptyPlaceholder,
    FormEdit,
    FormEmbed,
  },

  data() {
    return {
      curItem: null,
      isEditing: false,
      isFormVisible: false,
      isEmbedVisible: false,
    };
  },

  methods: {
    fetchForms() {
      this.$api.getForms();
    },

    showEditForm(data) {
      this.curItem = data;
      this.isFormVisible = true;
      this.isEditing = true;
    },

    showNewForm() {
      this.curItem = {};
      this.isFormVisible = true;
      this.isEditing = false;
    },

    showEmbed(data) {
      this.curItem = data;
      this.isEmbedVisible = true;
    },

    formFinished() {
      this.fetchForms();
    },

    deleteForm(f) {
      this.$utils.confirm(
        this.$t('forms.confirmDelete'),
        () => {
          this.$api.deleteForm(f.id).then(() => {
            this.fetchForms();
            this.$utils.toast(this.$t('globals.messages.deleted', { name: f.name }));
          });
        },
      );
    },
  },

  computed: {
    ...mapState(['forms', 'loading']),
  },

  created() {
    this.$root.$on('page.refresh', this.fetchForms);
  },

  destroyed() {
    this.$root.$off('page.refresh', this.fetchForms);
  },

  mounted() {
    this.fetchForms();
  },
});
</script>
