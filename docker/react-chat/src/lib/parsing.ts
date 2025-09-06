const tagRegex = /(\/\*(\w+)\*\/)\n?[\s\S]*?(?=\/\*|$)/g;

export const parseStructuredThoughts = (text: string) => {
  const parts: any[] = [];
  let match;
  tagRegex.lastIndex = 0; // Reset regex state

  while ((match = tagRegex.exec(text)) !== null) {
    const tagName = match[2].toLowerCase();
    const content = match[3].trim();
    if (content) {
      parts.push({ [tagName]: { content } });
    }
  }
  return parts;
};
