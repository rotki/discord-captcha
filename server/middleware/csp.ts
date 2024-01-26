const SELF = '\'self\'';
const UNSAFE_INLINE = '\'unsafe-inline\'';
const UNSAFE_EVAL = '\'unsafe-eval\'';
const NONE = '\'none\'';

const ContentPolicy = {
  BASE_URI: 'base-uri',
  BLOCK_ALL_MIXED_CONTENT: 'block-all-mixed-content',
  CHILD_SRC: 'child-src',
  CONNECT_SRC: 'connect-src',
  DEFAULT_SRC: 'default-src',
  FONT_SRC: 'font-src',
  FORM_ACTION: 'form-action',
  FRAME_ANCESTORS: 'frame-ancestors',
  FRAME_SRC: 'frame-src',
  IMG_SRC: 'img-src',
  OBJECT_SRC: 'object-src',
  SCRIPT_SRC: 'script-src',
  STYLE_SRC: 'style-src',
  WORKER_SRC: 'worker-src',
} as const;

type ContentPolicy = (typeof ContentPolicy)[keyof typeof ContentPolicy];

const policy: Record<ContentPolicy, string[]> = {
  [ContentPolicy.BASE_URI]: [SELF],
  [ContentPolicy.BLOCK_ALL_MIXED_CONTENT]: [],
  [ContentPolicy.CHILD_SRC]: [NONE],
  [ContentPolicy.CONNECT_SRC]: [SELF],
  [ContentPolicy.DEFAULT_SRC]: [SELF],
  [ContentPolicy.FONT_SRC]: [SELF, 'fonts.gstatic.com'],
  [ContentPolicy.FORM_ACTION]: [SELF],
  [ContentPolicy.FRAME_ANCESTORS]: [SELF],
  [ContentPolicy.FRAME_SRC]: [
    '*.recaptcha.net',
    'recaptcha.net',
    'https://www.google.com/recaptcha/',
    'https://recaptcha.google.com',
  ],
  [ContentPolicy.IMG_SRC]: [
    SELF,
    UNSAFE_INLINE,
    'data:',
    'www.gstatic.com/recaptcha',
  ],
  [ContentPolicy.OBJECT_SRC]: [NONE],
  [ContentPolicy.SCRIPT_SRC]: [
    SELF,
    UNSAFE_INLINE,
    UNSAFE_EVAL,
    'https://www.recaptcha.net',
    'https://recaptcha.net',
    'https://www.gstatic.com/recaptcha/',
    'https://www.gstatic.cn/recaptcha/',
    'https://www.google.com/recaptcha/',
  ],
  [ContentPolicy.STYLE_SRC]: [SELF, UNSAFE_INLINE, 'fonts.googleapis.com'],
  [ContentPolicy.WORKER_SRC]: [SELF, 'www.recaptcha.net'],
};

function getCSP() {
  let csp = '';
  const finalPolicy = { ...policy };
  for (const policyKey in finalPolicy)
    csp += `${policyKey} ${finalPolicy[policyKey as ContentPolicy].join(' ')};`;

  return csp;
}

const defaultCSP = getCSP();

export default defineEventHandler((event) => {
  event.node.res.setHeader('content-security-policy', defaultCSP);
});
