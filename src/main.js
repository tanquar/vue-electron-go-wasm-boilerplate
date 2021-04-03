import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

const go = new global.Go() // Loaded in public/index.html
const WASM_URL = '/wasm/main.wasm'

WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject).then(
  function (obj) {
    const wasm = obj.instance
    go.run(wasm)
  }
)

Vue.mixin({
  methods: {
    run: function (program, data) {
      return global.run(program, data)
    },
  },
})

Vue.config.productionTip = false

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount('#app')
