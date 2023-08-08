const SELF = "'self'";
const UNSAFE_INLINE = "'unsafe-inline'";
const UNSAFE_EVAL = "'unsafe-eval'";
const NONE = "'none'";

const ContentPolicy = {
  FRAME_ANCESTORS: 'frame-ancestors',
  BLOCK_ALL_MIXED_CONTENT: 'block-all-mixed-content',
  DEFAULT_SRC: 'default-src',
  SCRIPT_SRC: 'script-src',
  STYLE_SRC: 'style-src',
  OBJECT_SRC: 'object-src',
  FRAME_SRC: 'frame-src',
  CHILD_SRC: 'child-src',
  IMG_SRC: 'img-src',
  CONNECT_SRC: 'connect-src',
  BASE_URI: 'base-uri',
  FORM_ACTION: 'form-action',
  WORKER_SRC: 'worker-src',
  FONT_SRC: 'font-src',
} as const;

type ContentPolicy = (typeof ContentPolicy)[keyof typeof ContentPolicy];

const policy: Record<ContentPolicy, string[]> = {
  [ContentPolicy.FRAME_ANCESTORS]: [SELF],
  [ContentPolicy.BLOCK_ALL_MIXED_CONTENT]: [],
  [ContentPolicy.DEFAULT_SRC]: [SELF],
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
  [ContentPolicy.OBJECT_SRC]: [NONE],
  [ContentPolicy.FRAME_SRC]: [
    '*.recaptcha.net',
    'recaptcha.net',
    'https://www.google.com/recaptcha/',
    'https://recaptcha.google.com',
  ],
  [ContentPolicy.CHILD_SRC]: [NONE],
  [ContentPolicy.IMG_SRC]: [
    SELF,
    UNSAFE_INLINE,
    'data:',
    'www.gstatic.com/recaptcha',
  ],
  [ContentPolicy.CONNECT_SRC]: [SELF],
  [ContentPolicy.BASE_URI]: [SELF],
  [ContentPolicy.FORM_ACTION]: [SELF],
  [ContentPolicy.WORKER_SRC]: [SELF, 'www.recaptcha.net'],
  [ContentPolicy.FONT_SRC]: [SELF, 'fonts.gstatic.com'],
};

function getCSP() {
  let csp = '';
  const finalPolicy = { ...policy };
  for (const policyKey in finalPolicy) {
    csp += `${policyKey} ${finalPolicy[policyKey as ContentPolicy].join(' ')};`;
  }
  return csp;
}

const defaultCSP = getCSP();

export default defineEventHandler((event) => {
  event.node.res.setHeader('content-security-policy', defaultCSP);
});
