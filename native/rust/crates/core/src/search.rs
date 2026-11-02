//! Tool search. The algorithm is specified by `/protocol/fixtures/search.cases.json`
//! and mirrored in `src/core/search/search.ts`; both are tested against it.
//!
//! Rules: the query is split into tokens; every token must match at least one
//! field of the tool (AND). A token scores the best of:
//!   name word exact 100 | keyword exact 90 | name word prefix 80 |
//!   keyword prefix 70   | category exact 40 | description word prefix 30
//! A tool's score is the sum over tokens. Results are sorted by score
//! (descending) and ties keep the input order. An empty query returns the
//! input order.

use serde::Deserialize;

#[derive(Debug, Clone, Deserialize)]
pub struct Tool {
    pub id: String,
    pub name: String,
    #[serde(default)]
    pub description: String,
    #[serde(default)]
    pub category: String,
    #[serde(default)]
    pub keywords: Vec<String>,
}

/// Lowercase and split on every non-alphanumeric character.
pub fn words(s: &str) -> Vec<String> {
    s.split(|c: char| !c.is_alphanumeric())
        .filter(|w| !w.is_empty())
        .map(|w| w.to_lowercase())
        .collect()
}

fn best_in(token: &str, words: &[String], exact: u32, prefix: u32) -> u32 {
    let mut best = 0;
    for w in words {
        if w == token {
            return exact;
        }
        if w.starts_with(token) {
            best = best.max(prefix);
        }
    }
    best
}

struct Prepared {
    name: Vec<String>,
    keywords: Vec<String>,
    category: Vec<String>,
    description: Vec<String>,
}

impl Prepared {
    fn new(t: &Tool) -> Self {
        Self {
            name: words(&t.name),
            keywords: t.keywords.iter().flat_map(|k| words(k)).collect(),
            category: words(&t.category),
            description: words(&t.description),
        }
    }

    fn score_token(&self, token: &str) -> u32 {
        let mut best = best_in(token, &self.name, 100, 80);
        best = best.max(best_in(token, &self.keywords, 90, 70));
        if self.category.iter().any(|w| w == token) {
            best = best.max(40);
        }
        if self.description.iter().any(|w| w.starts_with(token)) {
            best = best.max(30);
        }
        best
    }
}

/// Returns the ids of matching tools, best first. `limit == 0` means no limit.
pub fn search(tools: &[Tool], query: &str, limit: usize) -> Vec<String> {
    let tokens = words(query);
    let mut scored: Vec<(u32, usize)> = Vec::new();
    for (i, t) in tools.iter().enumerate() {
        if tokens.is_empty() {
            scored.push((0, i));
            continue;
        }
        let p = Prepared::new(t);
        let mut total = 0;
        let mut all = true;
        for tok in &tokens {
            let s = p.score_token(tok);
            if s == 0 {
                all = false;
                break;
            }
            total += s;
        }
        if all {
            scored.push((total, i));
        }
    }
    // Stable sort: ties keep input order.
    scored.sort_by_key(|&(score, _)| std::cmp::Reverse(score));
    let it = scored.into_iter().map(|(_, i)| tools[i].id.clone());
    if limit > 0 {
        it.take(limit).collect()
    } else {
        it.collect()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use serde::Deserialize;
    use std::path::Path;

    #[derive(Deserialize)]
    struct Case {
        query: String,
        #[serde(default)]
        limit: usize,
        expected: Vec<String>,
    }
    #[derive(Deserialize)]
    struct File {
        tools: Vec<Tool>,
        cases: Vec<Case>,
    }

    #[test]
    fn matches_shared_fixture() {
        let p = Path::new(env!("CARGO_MANIFEST_DIR"))
            .join("../../../../protocol/fixtures/search.cases.json");
        let f: File = serde_json::from_str(&std::fs::read_to_string(p).unwrap()).unwrap();
        for c in &f.cases {
            assert_eq!(
                search(&f.tools, &c.query, c.limit),
                c.expected,
                "query {:?}",
                c.query
            );
        }
    }

    #[test]
    fn splits_on_punctuation() {
        assert_eq!(
            words("images/videos, JSON."),
            vec!["images", "videos", "json"]
        );
    }
}
