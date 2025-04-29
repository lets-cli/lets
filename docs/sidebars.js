module.exports = {
  mySidebar: [
    {
      type: 'category',
      label: 'Introduction',
      collapsed: false,
      items: [
        {
          type: 'doc',
          id: 'what_is_lets',
        },
        {
          type: 'doc',
          id: 'installation',
        },
        {
          type: 'doc',
          id: 'quick_start',
        },
        {
          type: 'doc',
          id: 'completion',
        },
      ],
    },
    {
      type: 'category',
      label: 'Usage',
      items: [
        {
          type: 'doc',
          id: 'basic_usage',
        },
        {
          type: 'doc',
          id: 'advanced_usage',
        },
      ],
    },
    'config',
    {
      type: 'category',
      label: 'API Reference',
      items: [
        {
          type: 'doc',
          id: 'cli',
        },
        {
          type: 'doc',
          id: 'env',
        },
      ],
    },
    {
      type: 'category',
      label: 'Examples',
      items: [
        {
          type: 'doc',
          id: 'examples',
        },
        {
          type: 'doc',
          id: 'example_js',
        },
      ],
    },
    'best_practices',
    'changelog',
    'ide_support',

    {
      type: 'category',
      label: 'Development',
      items: [
        {
          type: 'doc',
          id: 'architecture',
        },
        {
          type: 'doc',
          id: 'development',
        },
        {
          type: 'doc',
          id: 'contribute',
        },
      ],
    },
  ],
};
