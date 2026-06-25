package tool

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"net/smtp"
)

func Send_email(db *sql.DB, ip string, recipient string, title string, body string) error {
    rows := Query_DB(
        db,
        `select name, data from other where name in ("smtp_email", "smtp_pass", "smtp_server", "smtp_port", "smtp_security")`,
    )
    defer rows.Close()

    smtp_email := ""
    smtp_pass := ""
    smtp_server := ""
    smtp_security := ""
    smtp_port := ""

    for rows.Next() {
        var name, data string
        
        if err := rows.Scan(&name, &data); err != nil {
            return fmt.Errorf("failed to scan row: %v", err)
        }

        switch name {
        case "smtp_email":
            smtp_email = data
        case "smtp_pass":
            smtp_pass = data
        case "smtp_server":
            smtp_server = data
        case "smtp_security":
            smtp_security = data
        case "smtp_port":
            smtp_port = data
        }
    }

    if smtp_email == "" || smtp_pass == "" || smtp_server == "" || smtp_port == "" {
        return fmt.Errorf("smtp configuration is incomplete")
    }

    smtp_address := smtp_server + ":" + smtp_port
    var client *smtp.Client

    switch smtp_security {
    case "plain":
        conn, err := smtp.Dial(smtp_address)
        if err != nil {
            return fmt.Errorf("failed to connect to smtp server: %v", err)
        }
        client = conn
    case "starttls":
        conn, err := smtp.Dial(smtp_address)
        if err != nil {
            return fmt.Errorf("failed to connect to smtp server: %v", err)
        }
        if err := conn.StartTLS(&tls.Config{ServerName: smtp_server}); err != nil {
            return fmt.Errorf("failed to start tls: %v", err)
        }
        client = conn
    default:
        tls_conn, err := tls.Dial("tcp", smtp_address, &tls.Config{ServerName: smtp_server})
        if err != nil {
            return fmt.Errorf("failed to establish ssl connection: %v", err)
        }
        client, err = smtp.NewClient(tls_conn, smtp_server)
        if err != nil {
            return fmt.Errorf("failed to create smtp client: %v", err)
        }
    }

    defer client.Quit()

    auth := smtp.PlainAuth("", smtp_email, smtp_pass, smtp_server)
    if err := client.Auth(auth); err != nil {
        return fmt.Errorf("smtp authentication failed: %v", err)
    }

    if err := client.Mail(smtp_email); err != nil {
        return fmt.Errorf("failed to set sender: %v", err)
    }
    if err := client.Rcpt(recipient); err != nil {
        return fmt.Errorf("failed to set recipient: %v", err)
    }

    writer, err := client.Data()
    if err != nil {
        return fmt.Errorf("failed to send email data: %v", err)
    }

    domain := Get_domain(db, false)
    wiki_name := Get_wiki_set(db, ip, "")[0]

    message := fmt.Sprintf("from: %s <noreply@%s>\r\nto: %s\r\nsubject: %s\r\n\r\n%s", wiki_name, domain, recipient, title, body)

    _, err = writer.Write([]byte(message))
    if err != nil {
        return fmt.Errorf("failed to write email content: %v", err)
    }

    err = writer.Close()
    if err != nil {
        return fmt.Errorf("failed to finalize email send: %v", err)
    }

    return nil
}