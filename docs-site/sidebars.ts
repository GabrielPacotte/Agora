import type { SidebarsConfig } from '@docusaurus/plugin-content-docs';
import apiSidebar from './docs/api/sidebar';

const sidebars: SidebarsConfig = {
  codeSidebar: [
    {
      type: 'category',
      label: 'Code Docs',
      items: [
        'intro',
        {
          type: 'category',
          label: 'Architecture',
          items: [
            'architecture/project-structure',
            'architecture/domain',
            'architecture/sql',
            'architecture/http',
            'architecture/testing'
          ],
        },
      ],
    },
  ],
  apiSidebar
};

export default sidebars;
