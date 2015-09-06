/**
 * Wrap React app into html markup
 *
 * @param   Object   props  request options
 */
export default function(props) {
  return `<html lang="${props.lang || 'en'}" style="font-size: ${(props.fontSize || 100)}px">
    <head>
      <meta charset="UTF-8" />
      <link rel="stylesheet" href="/static/build/bundle.css"/>
      <link rel="icon" type="image/vnd.microsoft.icon" href="${require('#c/app/favicon.ico')}" />
      <title>${props.head.title}</title>
      ${props.head.meta}
      ${props.head.link}
    </head>
    <body>
      <div id="app">${props.app}</div>
      <script async src="/static/build/bundle.js"></script>
    </body>
  </html>`;
}
