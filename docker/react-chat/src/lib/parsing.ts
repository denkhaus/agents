const tagRegex = /(\/\*(\w+)\*\/)\n?[\s\S]*?(?=\/\*|$)/g;

export const parseStructuredThoughts = (text: string) => {
  const parts: any[] = [];
  let match;
  tagRegex.lastIndex = 0; // Reset regex state

  console.log("Parsing text:", text); // Add this line

  while ((match = tagRegex.exec(text)) !== null) {
    console.log("Match found:", match); // Add this line
    const tagName = match[2].toLowerCase();
    const content = match[3] ? match[3].trim() : "";
    if (content) {
      parts.push({ [tagName]: { content } });
    }
  }
  return parts;
};
