import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    blocks: [],
    totalBlocks: 0
  },
  mutations: {
    SET_BLOCKS (state, blocks) {
      state.blocks = blocks
    },
    SET_STATS (state, stats) {
      state.totalBlocks = stats.bestHeight
    }
  },
  actions: {
    loadBlocks ({ commit }, { page, rowsPerPage }) {
      return new Promise((resolve, reject) => {
        console.log(page, rowsPerPage)

        axios.get('/api/blocks', {
          params: {
            page: page - 1,
            limit: rowsPerPage
          }
        }).then(r => {
          commit('SET_BLOCKS', r.data.blocks)
          commit('SET_STATS', r.data.stats)
          resolve()
        })
      })
    }
  }
})
