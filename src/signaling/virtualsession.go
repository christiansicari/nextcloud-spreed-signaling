/**
 * Standalone signaling server for the Nextcloud Spreed app.
 * Copyright (C) 2019 struktur AG
 *
 * @author Joachim Bauch <bauch@struktur.de>
 *
 * @license GNU AGPL version 3 or any later version
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package signaling

import (
	"encoding/json"
	"net/url"
	"sync/atomic"
	"time"
	"unsafe"
)

type VirtualSession struct {
	hub       *Hub
	session   *ClientSession
	privateId string
	publicId  string
	data      *SessionIdData
	room      unsafe.Pointer

	sessionId string
	userId    string
	userData  *json.RawMessage
}

func GetVirtualSessionId(session *ClientSession, sessionId string) string {
	return session.PublicId() + "|" + sessionId
}

func NewVirtualSession(session *ClientSession, privateId string, publicId string, data *SessionIdData, msg *AddSessionInternalClientMessage) *VirtualSession {
	return &VirtualSession{
		session:   session,
		privateId: privateId,
		publicId:  publicId,
		data:      data,

		sessionId: msg.SessionId,
		userId:    msg.UserId,
		userData:  msg.User,
	}
}

func (s *VirtualSession) PrivateId() string {
	return s.privateId
}

func (s *VirtualSession) PublicId() string {
	return s.publicId
}

func (s *VirtualSession) ClientType() string {
	return HelloClientTypeVirtual
}

func (s *VirtualSession) Data() *SessionIdData {
	return s.data
}

func (s *VirtualSession) Backend() *Backend {
	return s.session.Backend()
}

func (s *VirtualSession) BackendUrl() string {
	return s.session.BackendUrl()
}

func (s *VirtualSession) ParsedBackendUrl() *url.URL {
	return s.session.ParsedBackendUrl()
}

func (s *VirtualSession) UserId() string {
	return s.userId
}

func (s *VirtualSession) UserData() *json.RawMessage {
	return s.userData
}

func (s *VirtualSession) SetRoom(room *Room) {
	atomic.StorePointer(&s.room, unsafe.Pointer(room))
}

func (s *VirtualSession) GetRoom() *Room {
	return (*Room)(atomic.LoadPointer(&s.room))
}

func (s *VirtualSession) LeaveRoom(notify bool) *Room {
	room := s.GetRoom()
	if room == nil {
		return nil
	}

	s.SetRoom(nil)
	room.RemoveSession(s)
	return room
}

func (s *VirtualSession) IsExpired(now time.Time) bool {
	return false
}

func (s *VirtualSession) Close() {
	s.session.RemoveVirtualSession(s)
	s.session.hub.removeSession(s)
}

func (s *VirtualSession) HasPermission(permission Permission) bool {
	return true
}

func (s *VirtualSession) Session() *ClientSession {
	return s.session
}

func (s *VirtualSession) SessionId() string {
	return s.sessionId
}
