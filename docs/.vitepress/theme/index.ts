import Theme from 'vitepress/theme'

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'


export default {
  ...Theme,
   enhanceApp: async({app, router, siteData}) => {
        app.use(ElementPlus)
   }
}