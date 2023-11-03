// emailHandler.js

document.addEventListener("DOMContentLoaded", function() {
    const reminderForm = document.querySelector(".reminder-form");

    reminderForm.addEventListener("submit", function(event) {
        event.preventDefault();
        
        // Get the form values using the placeholders from the provided index.html
        const recipientEmail = document.getElementById("recipientEmail").value;
        const senderName = document.getElementById("senderName").value;
        const reminderDateTime = document.getElementById("reminderDateTime").value;
        const reminderMessage = document.getElementById("reminderMessage").value;
        
        // Create the .ics content
        const icsContent = createIcsContent(reminderDateTime, reminderMessage);
        
        // Send the email data to the backend
        sendEmailDataToBackend(recipientEmail, senderName, icsContent, reminderMessage);
    });
});

function createIcsContent(dateTime, message) {
    const startTime = new Date(dateTime).toISOString().replace(/-|:|\.\d\d\d/g, "");
    const endTime = new Date(new Date(dateTime).setHours(new Date(dateTime).getHours() + 1)).toISOString().replace(/-|:|\.\d\d\d/g, "");

    return `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//hacksw/handcal//NONSGML v1.0//EN
BEGIN:VEVENT
UID:${new Date().getTime()}@keepinmind.com
DTSTAMP:${new Date().toISOString().replace(/-|:|\.\d\d\d/g, "")}
DTSTART:${startTime}
DTEND:${endTime}
SUMMARY:${message}
BEGIN:VALARM
TRIGGER:-PT15M
ACTION:DISPLAY
DESCRIPTION:Reminder
END:VALARM
END:VEVENT
END:VCALENDAR`;
}

function sendEmailDataToBackend(email, sender, icsContent, message) {
    const payload = {
        recipientEmail: email,
        sentBy: sender,
        ics: icsContent,
        reminderMessage: message
    };

    // Making a fetch call to the backend endpoint
    fetch("http://localhost:8080/sendReminder", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            alert("Email sent successfully!");
        } else {
            alert("Failed to send email. Please try again.");
        }
    })
    .catch(error => {
        console.error("Error sending email:", error);
    });
}
