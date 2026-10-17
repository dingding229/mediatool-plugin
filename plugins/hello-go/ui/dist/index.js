import { defineComponent, h, ref, components } from '@mediatool/plugin-ui';
export default defineComponent({props:{instanceId:String,api:Object},setup(props){const count=ref(0);return()=>h('section',{'aria-label':'Hello ESM'},[h('p','内置 Go 插件 · 共享 Vue 运行时'),h(components.ElButton,{onClick:()=>count.value++},()=>`交互 ${count.value}`)])}});
