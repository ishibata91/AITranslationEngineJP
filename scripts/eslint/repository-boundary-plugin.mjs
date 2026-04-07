import path from "node:path";

const SOURCE_LAYERS = ["ui", "application", "gateway", "shared"];
const ALLOWED_IMPORTS = {
  ui: new Set(["ui", "application", "shared"]),
  application: new Set(["application", "shared"]),
  gateway: new Set(["gateway", "application", "shared"]),
  shared: new Set(["shared"])
};
const PUBLIC_ROOT_SPEC_MATCHERS = [
  {
    layer: "ui",
    match(segments) {
      if (segments[0] === "app-shell") {
        return {
          rootId: "ui:app-shell",
          rootSegmentLength: 1
        };
      }

      if (segments[0] === "screens" && segments[1]) {
        return {
          rootId: `ui:screens/${segments[1]}`,
          rootSegmentLength: 2
        };
      }

      if (segments[0] === "views") {
        return {
          rootId: "ui:views",
          rootSegmentLength: 1
        };
      }

      if (segments[0] === "stores") {
        return {
          rootId: "ui:stores",
          rootSegmentLength: 1
        };
      }

      return null;
    }
  },
  {
    layer: "application",
    match(segments) {
      if (segments[0] === "bootstrap") {
        return {
          rootId: "application:bootstrap",
          rootSegmentLength: 1
        };
      }

      if (segments[0] === "usecases") {
        return {
          rootId: "application:usecases",
          rootSegmentLength: 1
        };
      }

      if (segments[0] === "ports" && segments[1] === "input") {
        return {
          rootId: "application:ports/input",
          rootSegmentLength: 2
        };
      }

      if (segments[0] === "ports" && segments[1] === "gateway") {
        return {
          rootId: "application:ports/gateway",
          rootSegmentLength: 2
        };
      }

      return null;
    }
  },
  {
    layer: "gateway",
    match(segments) {
      if (segments[0] === "wails") {
        return {
          rootId: "gateway:wails",
          rootSegmentLength: 1
        };
      }

      return null;
    }
  },
  {
    layer: "shared",
    match(segments) {
      if (segments[0] === "contracts") {
        return {
          rootId: "shared:contracts",
          rootSegmentLength: 1
        };
      }

      return null;
    }
  }
];

function normalizePath(value) {
  return value.replaceAll("\\", "/");
}

function isTestLikePath(value) {
  return /(^|\/)[^/]+\.test(?:\.[^/]+)?$/u.test(value);
}

function detectSourceLayer(filename) {
  const normalized = normalizePath(filename);

  for (const layer of SOURCE_LAYERS) {
    if (normalized.includes(`/src/${layer}/`)) {
      return layer;
    }
  }

  return null;
}

function detectResolvedTargetType(resolvedPath) {
  const normalized = normalizePath(resolvedPath);

  if (
    normalized.includes("/src/test/") ||
    normalized.includes("/fixtures/") ||
    normalized.includes("/generated/") ||
    isTestLikePath(normalized)
  ) {
    return "reverse-flow";
  }

  for (const layer of SOURCE_LAYERS) {
    if (normalized.includes(`/src/${layer}/`)) {
      return layer;
    }
  }

  return null;
}

function detectTargetType(filename, specifier) {
  if (
    specifier === "wailsjs" ||
    specifier.startsWith("wailsjs/") ||
    specifier.includes("/wailsjs/")
  ) {
    return "wails";
  }

  for (const layer of SOURCE_LAYERS) {
    if (specifier === `@${layer}` || specifier.startsWith(`@${layer}/`)) {
      return layer;
    }
  }

  if (specifier.startsWith(".")) {
    return detectResolvedTargetType(path.resolve(path.dirname(filename), specifier));
  }

  return null;
}

function isAliasSpecifier(specifier) {
  for (const layer of SOURCE_LAYERS) {
    if (specifier === `@${layer}` || specifier.startsWith(`@${layer}/`)) {
      return true;
    }
  }

  return false;
}

function splitAliasSpecifier(specifier) {
  if (!isAliasSpecifier(specifier)) {
    return null;
  }

  const [layerAlias, ...segments] = specifier.split("/");
  return { layer: layerAlias.slice(1), segments };
}

function stripSrcPrefixSegments(normalizedAbsolutePath) {
  const segments = normalizedAbsolutePath.split("/");
  const srcIndex = segments.lastIndexOf("src");

  if (srcIndex < 0 || srcIndex + 1 >= segments.length) {
    return null;
  }

  return {
    layer: segments[srcIndex + 1],
    segments: segments.slice(srcIndex + 2)
  };
}

function detectPublicRoot(layer, segments) {
  for (const matcher of PUBLIC_ROOT_SPEC_MATCHERS) {
    if (matcher.layer !== layer) {
      continue;
    }

    const result = matcher.match(segments);

    if (result !== null) {
      return {
        layer,
        rootId: result.rootId,
        remainderSegments: segments.slice(result.rootSegmentLength)
      };
    }
  }

  return null;
}

function detectSourcePublicRoot(filename) {
  const normalized = normalizePath(filename);
  const parsed = stripSrcPrefixSegments(normalized);

  if (parsed === null) {
    return null;
  }

  return detectPublicRoot(parsed.layer, parsed.segments);
}

function detectTargetPublicRoot(filename, specifier) {
  if (isAliasSpecifier(specifier)) {
    const parsed = splitAliasSpecifier(specifier);

    if (parsed === null) {
      return null;
    }

    return detectPublicRoot(parsed.layer, parsed.segments);
  }

  if (!specifier.startsWith(".")) {
    return null;
  }

  const resolvedPath = normalizePath(path.resolve(path.dirname(filename), specifier));
  const parsed = stripSrcPrefixSegments(resolvedPath);

  if (parsed === null) {
    return null;
  }

  return detectPublicRoot(parsed.layer, parsed.segments);
}

function isPublicEntrypointImport(targetPublicRoot) {
  const depth = targetPublicRoot.remainderSegments.length;

  if (depth === 0) {
    return true;
  }

  if (depth === 1) {
    return true;
  }

  return false;
}

function isReverseFlowTargetSpecifier(specifier) {
  if (specifier.includes("/fixtures/") || specifier.includes("/generated/")) {
    return true;
  }

  return isTestLikePath(normalizePath(specifier));
}

function isReverseFlowSourceExempt(filename) {
  const normalized = normalizePath(filename);
  return normalized.includes("/src/test/") || isTestLikePath(normalized);
}

function buildMessage(sourceLayer, targetType) {
  if (targetType === "wails") {
    return `${sourceLayer} code must not import Wails bindings directly. Go through gateway ports or gateway adapters instead.`;
  }

  if (targetType === "reverse-flow") {
    return `${sourceLayer} production code must not import test, fixture, or generated support files.`;
  }

  return `${sourceLayer} code must not import ${targetType} code directly.`;
}

function buildSameLayerInternalImportMessage(sourcePublicRoot, targetPublicRoot) {
  return `${sourcePublicRoot.rootId} must not import internal modules of ${targetPublicRoot.rootId}. Use the target root index or a direct child file.`;
}

const enforceLayerBoundariesRule = {
  meta: {
    type: "problem",
    docs: {
      description: "Enforce repository layer boundaries for frontend source code."
    },
    schema: [],
    messages: {
      forbiddenImport: "{{message}}"
    }
  },
  create(context) {
    const filename = context.filename ?? context.getFilename();
    const sourceLayer = detectSourceLayer(filename);
    const sourcePublicRoot = detectSourcePublicRoot(filename);

    if (sourceLayer === null) {
      return {};
    }

    return {
      ImportDeclaration(node) {
        const specifier = node.source.value;

        if (typeof specifier !== "string") {
          return;
        }

        const targetType = isReverseFlowTargetSpecifier(specifier)
          ? "reverse-flow"
          : detectTargetType(filename, specifier);

        if (targetType === null) {
          return;
        }

        if (targetType === "reverse-flow") {
          if (isReverseFlowSourceExempt(filename)) {
            return;
          }

          context.report({
            node: node.source,
            messageId: "forbiddenImport",
            data: {
              message: buildMessage(sourceLayer, targetType)
            }
          });
          return;
        }

        if (targetType === "wails" && sourceLayer === "gateway") {
          return;
        }

        if (!ALLOWED_IMPORTS[sourceLayer].has(targetType)) {
          context.report({
            node: node.source,
            messageId: "forbiddenImport",
            data: {
              message: buildMessage(sourceLayer, targetType)
            }
          });
          return;
        }

        if (targetType !== sourceLayer) {
          return;
        }

        if (sourcePublicRoot === null) {
          return;
        }

        const targetPublicRoot = detectTargetPublicRoot(filename, specifier);

        if (targetPublicRoot === null) {
          return;
        }

        if (targetPublicRoot.layer !== sourcePublicRoot.layer) {
          return;
        }

        if (targetPublicRoot.rootId === sourcePublicRoot.rootId) {
          return;
        }

        if (!isPublicEntrypointImport(targetPublicRoot)) {
          context.report({
            node: node.source,
            messageId: "forbiddenImport",
            data: {
              message: buildSameLayerInternalImportMessage(sourcePublicRoot, targetPublicRoot)
            }
          });
        }
      }
    };
  }
};

export const repositoryBoundaryPlugin = {
  meta: {
    name: "repository-boundary-plugin"
  },
  rules: {
    "enforce-layer-boundaries": enforceLayerBoundariesRule
  }
};
