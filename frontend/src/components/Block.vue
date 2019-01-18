<template>
  <div>
    <v-breadcrumbs
      divider=">"
      :items="breadcrumbs"
    >
    </v-breadcrumbs>

    <v-card v-if="header !== null">
      <v-card-title>
        <h4>Summary</h4>
      </v-card-title>
      <v-divider></v-divider>
      <v-layout row wrap>
        <v-flex sm12 md3 lg3 xl3>
          <v-list dense>
            <v-list-tile>
              <v-list-tile-content class="body-2">Height</v-list-tile-content>
              <v-list-tile-content class="align-end">
                {{ header.height }}
              </v-list-tile-content>
            </v-list-tile>
          </v-list>
          <v-list dense>
            <v-list-tile>
              <v-list-tile-content class="body-2">Version</v-list-tile-content>
              <v-list-tile-content class="align-end">
                {{ header.version }}
              </v-list-tile-content>
            </v-list-tile>
          </v-list>
          <v-list dense>
            <v-list-tile>
              <v-list-tile-content class="body-2">Flags</v-list-tile-content>
              <v-list-tile-content class="align-end">
                {{ header.flags }}
              </v-list-tile-content>
            </v-list-tile>
          </v-list>
        </v-flex>

        <v-flex sm12 md3 lg3 xl3>
          <v-list dense>
            <v-list-tile>
              <v-list-tile-content class="body-2">Bits</v-list-tile-content>
              <v-list-tile-content class="align-end">
                {{ header.bits }}
              </v-list-tile-content>
            </v-list-tile>
          </v-list>
          <v-list dense>
            <v-list-tile>
              <v-list-tile-content class="body-2">Nonce</v-list-tile-content>
              <v-list-tile-content class="align-end">
                {{ header.nonce }}
              </v-list-tile-content>
            </v-list-tile>
          </v-list>
          <v-list dense>
            <v-list-tile>
              <v-list-tile-content class="body-2">Time</v-list-tile-content>
              <v-list-tile-content class="align-end">
                {{ header.timestamp | moment('llll') }}
              </v-list-tile-content>
            </v-list-tile>
          </v-list>
        </v-flex>

        <v-flex sm12 md6 lg6 xl6>
          <v-list dense>
            <v-list-tile>
              <v-list-tile-content class="body-2">Block Hash</v-list-tile-content>
              <v-list-tile-content class="align-end">
                {{ header.hash }}
              </v-list-tile-content>
            </v-list-tile>
          </v-list>
          <v-list dense>
            <v-list-tile>
              <v-list-tile-content class="body-2">Prev Block</v-list-tile-content>
              <v-list-tile-content class="align-end">
                {{ header.hash_prev_block }}
              </v-list-tile-content>
            </v-list-tile>
          </v-list>
          <v-list dense>
            <v-list-tile>
              <v-list-tile-content class="body-2">Merkle Root</v-list-tile-content>
              <v-list-tile-content class="align-end">
                {{ header.hash_merkle_root }}
              </v-list-tile-content>
            </v-list-tile>
          </v-list>
        </v-flex>
      </v-layout>
    </v-card>

    <v-layout
      v-if="block !== null"
      column
      wrap
    >
      <v-flex xs12>
        <v-card
          v-for="tx in block.txs"
          :key="tx.hash"
          class="mt-3"
        >
          <v-card-title>
            <h4>{{ tx.hash }}</h4>
          </v-card-title>
          <v-divider></v-divider>
          <v-layout row fill-height class="pa-3">
            <v-flex md5 v-if="tx.inputs === null || tx.inputs.length === 0">
              <v-chip>
                Coinbase
              </v-chip>
            </v-flex>

            <v-flex md7>
              <v-chip
                v-for="output in tx.outputs"
                :key="output.script"
              >
                {{ output.value }}
              </v-chip>
            </v-flex>
          </v-layout>
        </v-card>
      </v-flex>
    </v-layout>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  props: {
    hash: String
  },
  data () {
    return {
      block: null,
      breadcrumbs: [
        {
          text: 'Home',
          disabled: false,
          to: '/'
        },
        {
          text: 'Block - ' + this.hash,
          disabled: true
        }
      ]
    }
  },
  mounted () {
    this.getBlock()
  },
  methods: {
    getBlock () {
      axios.get('/api/blocks/' + this.hash).then(res => {
        this.block = res.data.block
      })
    }
  },
  computed: {
    header () {
      return this.block == null ? null : this.block.header
    }
  }
}
</script>

<style>

</style>
