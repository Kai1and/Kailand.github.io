import { ImagePlus, MessageCircle, Send } from "lucide-react";
import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { chatApi } from "../api/chatApi.js";
import { useAuthStore } from "../store/authStore.js";

export default function ChatPage() {
  const { id } = useParams();
  const { user } = useAuthStore();
  const [conversations, setConversations] = useState([]);
  const [messages, setMessages] = useState([]);
  const [body, setBody] = useState("");
  const [attachment, setAttachment] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    chatApi.list().then(setConversations).catch((err) => setError(err.message));
  }, []);

  useEffect(() => {
    if (id) {
      chatApi.messages(id).then(setMessages).catch((err) => setError(err.message));
    }
  }, [id]);

  const send = async (event) => {
    event.preventDefault();
    if (!id || (!body.trim() && !attachment)) return;
    const message = await chatApi.send(id, { body: body.trim(), attachment_url: attachment });
    setMessages([...messages, message]);
    setBody("");
    setAttachment("");
    chatApi.list().then(setConversations).catch(() => {});
  };

  const pickAttachment = (event) => {
    const file = event.target.files?.[0];
    if (!file) return;
    const reader = new FileReader();
    reader.onload = () => setAttachment(reader.result);
    reader.readAsDataURL(file);
  };

  const activeConversation = conversations.find((conversation) => String(conversation.id) === String(id));

  return (
    <div className="chat-layout">
      <aside className="chat-list panel">
        <div className="panel-head">
          <h2>Сообщения</h2>
        </div>
        {error && <div className="notice">{error}</div>}
        {conversations.map((conversation) => (
          <Link className={`chat-preview ${String(conversation.id) === String(id) ? "active" : ""}`} to={`/chats/${conversation.id}`} key={conversation.id}>
            <img src={conversation.equipment?.image_url} alt="" />
            <strong>{conversation.equipment?.name ?? `Диалог #${conversation.id}`}</strong>
            {conversation.unread_count > 0 && <em className="unread-pill">{conversation.unread_count}</em>}
            <small>
              {user?.id === conversation.owner_id ? conversation.customer?.name : conversation.owner?.name}
            </small>
            <span>{conversation.last_message || "Нет сообщений"}</span>
          </Link>
        ))}
        {!conversations.length && !error && <div className="empty-state">Диалоги появятся после сообщения продавцу.</div>}
      </aside>

      <section className="chat-window panel">
        <div className="panel-head">
          <h2>{activeConversation?.equipment?.name ?? (id ? `Диалог #${id}` : "Выберите диалог")}</h2>
          {activeConversation && (
            <span className="chat-person">
              <MessageCircle size={16} />
              {user?.id === activeConversation.owner_id ? activeConversation.customer?.name : activeConversation.owner?.name}
            </span>
          )}
        </div>
        <div className="messages">
          {messages.map((message) => (
            <div className={message.sender_id === user?.id ? "message mine" : "message"} key={message.id}>
              {message.attachment_url && <img className="message-attachment" src={message.attachment_url} alt="Вложение" />}
              {message.body && <span>{message.body}</span>}
            </div>
          ))}
          {id && !messages.length && <div className="empty-state">Сообщений пока нет.</div>}
        </div>
        {id && (
          <form className="message-form" onSubmit={send}>
            <input value={body} onChange={(e) => setBody(e.target.value)} placeholder="Введите сообщение" />
            <label className="icon-button" aria-label="Прикрепить фото">
              <ImagePlus size={18} />
              <input type="file" accept="image/*" onChange={pickAttachment} />
            </label>
            <button className="primary-button">
              <Send size={18} />
              Отправить
            </button>
            {attachment && <img className="attachment-preview" src={attachment} alt="Предпросмотр" />}
          </form>
        )}
      </section>
    </div>
  );
}
