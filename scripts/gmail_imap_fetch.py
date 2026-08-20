#!/usr/bin/env python3
"""gmail_imap_fetch.py -- on-demand, single-query Gmail read tool.

Uses the SAME Gmail App Password credential as emily-agent/gmail.go's SMTP
send path (GMAIL_SMTP_ADDRESS / GMAIL_SMTP_PASSWORD) -- app passwords work
for both SMTP and IMAP, they aren't scoped per-protocol. Deliberately a
plain script, not a service: 2026-08-20, founder real-time, "for now it
needs to be treated as a queryable source not something you are totally
chunking through yet" -- on-demand lookups only, no background sync/index
job. Also honors the founder's explicit privacy boundary from the same
conversation: this script prints message content to stdout for the
operator running it interactively -- nothing it does writes email content
into Apples, BACKLOG.md, CHANGELOG.md, or any other persistent/git-tracked
file. Keep it that way if this ever gets extended.

Credentials: run via emily's own env-sourcing convention --
    source /home/fatbaby/EMILY/var/emily-secrets.env
    python3 gmail_imap_fetch.py --search 'subject:"language spec"'
    python3 gmail_imap_fetch.py --latest 5
    python3 gmail_imap_fetch.py --uid 12345

Stdlib only (imaplib + email), no dependencies to install.
"""

import argparse
import email
import email.header
import imaplib
import os
import sys

IMAP_HOST = "imap.gmail.com"
IMAP_PORT = 993


def decode_header_value(raw):
    if raw is None:
        return ""
    parts = email.header.decode_header(raw)
    out = []
    for text, enc in parts:
        if isinstance(text, bytes):
            out.append(text.decode(enc or "utf-8", errors="replace"))
        else:
            out.append(text)
    return "".join(out)


def extract_body(msg):
    if msg.is_multipart():
        for part in msg.walk():
            ctype = part.get_content_type()
            disp = str(part.get("Content-Disposition") or "")
            if ctype == "text/plain" and "attachment" not in disp:
                payload = part.get_payload(decode=True)
                if payload:
                    charset = part.get_content_charset() or "utf-8"
                    return payload.decode(charset, errors="replace")
        return "(no text/plain part found)"
    else:
        payload = msg.get_payload(decode=True)
        if payload:
            charset = msg.get_content_charset() or "utf-8"
            return payload.decode(charset, errors="replace")
        return str(msg.get_payload())


def connect():
    address = os.environ.get("GMAIL_SMTP_ADDRESS")
    password = os.environ.get("GMAIL_SMTP_PASSWORD")
    if not address or not password:
        print("error: GMAIL_SMTP_ADDRESS / GMAIL_SMTP_PASSWORD not set in environment", file=sys.stderr)
        print("       source /home/fatbaby/EMILY/var/emily-secrets.env first", file=sys.stderr)
        sys.exit(1)
    conn = imaplib.IMAP4_SSL(IMAP_HOST, IMAP_PORT)
    conn.login(address, password)
    conn.select("INBOX")
    return conn


def print_message(conn, uid):
    typ, data = conn.fetch(uid, "(RFC822)")
    if typ != "OK" or not data or data[0] is None:
        print(f"could not fetch uid {uid!r}", file=sys.stderr)
        return
    raw = data[0][1]
    msg = email.message_from_bytes(raw)
    print("=" * 70)
    print(f"From:    {decode_header_value(msg.get('From'))}")
    print(f"Subject: {decode_header_value(msg.get('Subject'))}")
    print(f"Date:    {msg.get('Date')}")
    print("-" * 70)
    print(extract_body(msg))
    print("=" * 70)


def main():
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--search", help='IMAP SEARCH query, e.g. subject:"language spec" or from:someone@example.com')
    ap.add_argument("--latest", type=int, metavar="N", help="print the N most recent messages in the inbox")
    ap.add_argument("--uid", help="fetch one specific message by UID")
    args = ap.parse_args()

    conn = connect()
    try:
        if args.uid:
            print_message(conn, args.uid)
            return

        if args.search:
            # Gmail's IMAP extension accepts raw Gmail search syntax via X-GM-RAW.
            typ, data = conn.uid("search", None, "X-GM-RAW", f'"{args.search}"')
        else:
            typ, data = conn.uid("search", None, "ALL")

        if typ != "OK":
            print("search failed", file=sys.stderr)
            sys.exit(1)

        uids = data[0].split()
        if args.latest:
            uids = uids[-args.latest:]

        if not uids:
            print("no matching messages")
            return

        for uid in uids:
            print_message(conn, uid)
    finally:
        try:
            conn.close()
        except Exception:
            pass
        conn.logout()


if __name__ == "__main__":
    main()
