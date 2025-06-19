export type Edge = {
  id: string;
  source: string;
  target: string;
  // "step" or "smoothstep" etc.
  type?: string;
};
