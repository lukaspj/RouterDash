import { describe, it, expect } from "bun:test";
import { readFileSync } from "fs";
import { join } from "path";

const html = readFileSync(join(import.meta.dir, "..", "index.html"), "utf8");

// Extract the <script> block
const scriptMatch = html.match(/<script>([\s\S]*?)<\/script>/);
if (!scriptMatch) throw new Error("No <script> block found in index.html");

const jsSource = scriptMatch[1];

describe("embedded JavaScript", () => {
  it("parses without syntax errors", () => {
    // new Function parses the code body — throws SyntaxError if invalid
    // Wrap in a function so 'return' at top level doesn't break it
    expect(() => new Function(jsSource)).not.toThrow();
  });

  it("contains dashboard() function", () => {
    expect(jsSource).toMatch("function dashboard()");
  });

  it("has all required fetch methods", () => {
    const fns = [
      "fetchSystem", "fetchResources", "fetchHardware",
      "fetchInterfaces", "fetchDhcp", "fetchFirewall", "fetchWireless",
    ];
    for (const fn of fns) {
      expect(jsSource).toMatch(new RegExp(`async ${fn}\\(\\)`));
    }
  });

  it("has all required data properties in dashboard() return object", () => {
    const props = [
      "activeTab", "tabs", "refreshing", "connected", "lastRefresh",
      "system", "resources", "resourcesFetched", "interfaces",
      "dhcpLeases", "firewallRules", "wirelessInterfaces",
      "bridgePorts", "ethernetInterfaces", "hardware",
      "firewallLoading", "firewallError",
      "interfacesLoading", "interfacesError",
      "dhcpLoading", "dhcpError",
      "wirelessLoading", "wirelessError",
      "hardwareLoading", "hardwareError",
      "resourcesLoading", "resourcesError",
      "systemLoading", "systemError",
    ];
    for (const prop of props) {
      expect(jsSource).toMatch(new RegExp(`\\b${prop}\\b.*:`));
    }
  });

  it("fwChainFlow getter exists", () => {
    expect(jsSource).toMatch(/get fwChainFlow/);
  });

  it("has no orphaned catch blocks (stray brace would leave catch detached)", () => {
    // Quick structural check: count try vs catch
    const tryCount = (jsSource.match(/\btry\b\s*\{/g) || []).length;
    const catchCount = (jsSource.match(/\}\s*catch\b/g) || []).length;
    expect(catchCount).toBe(tryCount);
  });
});
