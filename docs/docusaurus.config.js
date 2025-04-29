module.exports = {
  title: 'Lets',
  tagline: 'CLI task runner for developers - a better alternative to make',
  url: 'https://lets-cli.org',
  baseUrl: '/',
  favicon: 'img/favicon.ico',
  organizationName: 'lets-cli', // Usually your GitHub org/user name.
  projectName: 'lets', // Usually your repo name.
  themeConfig: {
    prism: {
      theme: require('prism-react-renderer/themes/vsDark'),
    },
    navbar: {
      title: 'Lets',
      logo: {
        alt: 'Lets Logo',
        src: 'img/logo.png',
      },
      items: [
        {to: 'docs/quick_start', label: 'Docs', position: 'right'},
        {to: 'docs/changelog', label: 'Changelog', position: 'right'},
        {to: 'blog', label: 'Blog', position: 'right'},
        {
          href: 'https://github.com/lets-cli/lets',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {
              label: 'Quick start',
              to: 'docs/quick_start',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'Stack Overflow',
              href: 'https://stackoverflow.com/questions/tagged/lets-cli',
            },
          ],
        },
        {
          title: 'More',
          items: [
            {
              label: 'Blog',
              to: 'blog',
            },
            {
              label: 'GitHub',
              href: 'https://github.com/facebook/docusaurus',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Lets, Inc. Built with Docusaurus.`,
    },
    algolia: {
      appId: "B314NWJQO4",
      apiKey: "3103c243857b4a1debe49df0c8206704",
      indexName: "lets-cli",
    }
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          editUrl:
            'https://github.com/lets-cli/lets/edit/master/docs/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
        gtag: {
          trackingID: 'G-DLCLPWY8PL',
          anonymizeIP: true,
        },
      },
    ],
  ],
};
