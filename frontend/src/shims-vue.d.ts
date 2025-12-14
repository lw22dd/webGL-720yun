declare module "element-plus";
declare module 'element-plus/dist/locale/zh-cn.mjs';
/* eslint-disable */
declare module '*.vue' {
    import type { DefineComponent } from 'vue';
    const component: DefineComponent<{}, {}, any>;
    export default component;
}

declare module '*.css' {
    const content: any;
    export default content;
}
declare module '*.scss' {
    const content: any;
    export default content;
}