<template>
  <Block :hash="$route.params.blockHash" />
</template>

<script>
import Block from '../components/Block'

export default {
  components: {
    Block
  },
  head: {
    title () {
      return {
        inner: 'Block - ' + this.$route.params.blockHash
      }
    }
  }
}
</script>
