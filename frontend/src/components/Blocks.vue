<template>
    <v-layout
      text-xs-center
      wrap
      align-center
      justify-center
    >
      <template>
        <v-data-table
          :headers="headers"
          :items="blocks"
          :total-items="totalBlocks"
          :pagination.sync="pagination"
          :loading="loading"
          class="elevation-1"
        >
          <template slot="no-data">
            <v-alert :value="true" color="error" icon="warning">
              Sorry, nothing to display here :(
            </v-alert>
          </template>

          <template slot="items" slot-scope="props">
            <td class="text-xs-left"><router-link :to="{ name: 'block', params: { blockHash: props.item.header.hash }}">{{ props.item.header.height }}</router-link></td>
            <td class="text-xs-left">{{ props.item.header.timestamp | moment('from') }}</td>
            <td class="text-xs-left">{{ props.item.txs === null ? 0 : props.item.txs.length }}</td>
            <td class="text-xs-left">1000</td>
          </template>
        </v-data-table>
      </template>
    </v-layout>
</template>

<script>
import { mapState } from 'vuex'

export default {
  data () {
    return {
      loading: true,
      pagination: {
        rowsPerPage: 10
      },
      headers: [
        {
          text: 'Height',
          align: 'left',
          sortable: false,
          value: 'height'
        },
        {
          text: 'Age',
          align: 'left',
          sortable: false,
          value: 'age'
        },
        {
          text: 'Transactions',
          align: 'left',
          sortable: false,
          value: 'transactions'
        },
        {
          text: 'Size',
          align: 'left',
          sortable: false,
          value: 'size'
        }
      ]
    }
  },
  watch: {
    pagination: {
      handler () {
        this.getBlocks()
      },
      deep: true
    }
  },
  mounted () {
    this.getBlocks()
  },
  methods: {
    getBlocks () {
      this.loading = true

      this.$store.dispatch('loadBlocks', this.pagination).then(() => {
        this.loading = false
      })
    }
  },
  computed: mapState([
    'blocks',
    'totalBlocks'
  ])
}
</script>

<style>

</style>
