# HIS API Specification

## Overview

Hospital Information System (HIS) RESTful API with multi-tenant architecture.

| Property | Value |
|----------|-------|
| **Base URL** | `https://{subdomain}.yourdomain.com/api/v1` |
| **Content-Type** | `application/json` |
| **Authentication** | Bearer Token (JWT) |

## Multi-Tenant Architecture

All API endpoints require tenant identification via subdomain:
- Requests to `https://bangkok.his.com/api/v1/...` target the Bangkok hospital tenant
- Requests to `https://chiangmai.his.com/api/v1/...` target the Chiang Mai hospital tenant

---

## Response Format

### Success Response

```json
{
  "success": true,
  "message": "operation message",
  "data": { ... }
}
```

### Error Response

```json
{
  "success": false,
  "error": "error message"
}
```

---

## Endpoints

### Health Check

#### `GET /health`

Check API health status. **No authentication required.**

**Response:**
```json
{
  "status": "ok"
}
```

---

## Staff Endpoints

### Login

#### `POST /api/v1/staff/login`

Authenticate staff and receive JWT token.

**Authentication:** None (Public)  
**Tenant Required:** Yes

**Request Body:**
| Field | Type | Required | Validation | Description |
|-------|------|----------|------------|-------------|
| `username` | string | ✅ | min=5, max=100 | Staff username |
| `password` | string | ✅ | min=6 | Staff password |
| `hospital_id` | uint | ✅ | - | Hospital ID |

**Request Example:**
```json
{
  "username": "admin",
  "password": "password123",
  "hospital_id": 1
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "staff": {
      "id": 1,
      "username": "admin",
      "staff_code": "STAFF-BKGH0001-000001",
      "phone_number": "0812345678",
      "email": "admin@hospital.com",
      "first_name": "John",
      "last_name": "Doe",
      "hospital_id": 1,
      "is_admin": true
    }
  }
}
```

**Error Response (401):**
```json
{
  "success": false,
  "error": "invalid credentials"
}
```

---

### Logout

#### `POST /api/v1/staff/logout`

Logout staff (client should discard the token).

**Authentication:** Bearer Token  
**Tenant Required:** Yes

**Headers:**
```
Authorization: Bearer <token>
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "logout successful",
  "data": null
}
```

---

### Get All Staff

#### `GET /api/v1/staff/`

Retrieve all staff members.

**Authentication:** Bearer Token  
**Tenant Required:** Yes

**Success Response (200):**
```json
{
  "success": true,
  "message": "success",
  "data": [
    {
      "ID": 1,
      "CreatedAt": "2024-01-01T00:00:00Z",
      "UpdatedAt": "2024-01-01T00:00:00Z",
      "DeletedAt": null,
      "username": "admin",
      "staff_code": "STAFF-BKGH0001-000001",
      "phone_number": "0812345678",
      "email": "admin@hospital.com",
      "first_name": "John",
      "last_name": "Doe",
      "hospital_id": 1,
      "is_admin": true
    }
  ]
}
```

---

### Get Staff by ID

#### `GET /api/v1/staff/:id`

Retrieve a specific staff member by ID.

**Authentication:** Bearer Token  
**Tenant Required:** Yes

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | uint | Staff ID |

**Success Response (200):**
```json
{
  "success": true,
  "message": "success",
  "data": {
    "ID": 1,
    "CreatedAt": "2024-01-01T00:00:00Z",
    "UpdatedAt": "2024-01-01T00:00:00Z",
    "DeletedAt": null,
    "username": "admin",
    "staff_code": "STAFF-BKGH0001-000001",
    "phone_number": "0812345678",
    "email": "admin@hospital.com",
    "first_name": "John",
    "last_name": "Doe",
    "hospital_id": 1,
    "is_admin": true
  }
}
```

**Error Response (404):**
```json
{
  "success": false,
  "error": "staff not found"
}
```

---

### Create Staff

#### `POST /api/v1/staff/create`

Create a new staff member. **Admin only.**

**Authentication:** Bearer Token (Admin)  
**Tenant Required:** Yes

**Request Body:**
| Field | Type | Required | Validation | Description |
|-------|------|----------|------------|-------------|
| `username` | string | ✅ | min=5, max=100 | Unique username |
| `password` | string | ✅ | min=6 | Password |
| `staff_code` | string | ❌ | - | Auto-generated if empty |
| `phone_number` | string | ✅ | - | Contact phone |
| `email` | string | ✅ | email format | Email address |
| `first_name` | string | ✅ | max=255 | First name |
| `last_name` | string | ✅ | max=255 | Last name |
| `hospital_id` | uint | ✅ | - | Hospital ID |
| `is_admin` | bool | ❌ | default=false | Admin privileges |

**Request Example:**
```json
{
  "username": "nurse001",
  "password": "password123",
  "phone_number": "0812345678",
  "email": "nurse@hospital.com",
  "first_name": "Jane",
  "last_name": "Smith",
  "hospital_id": 1,
  "is_admin": false
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "staff created successfully",
  "data": {
    "ID": 2,
    "CreatedAt": "2024-01-01T00:00:00Z",
    "UpdatedAt": "2024-01-01T00:00:00Z",
    "DeletedAt": null,
    "username": "nurse001",
    "staff_code": "STAFF-BKGH0001-000002",
    "phone_number": "0812345678",
    "email": "nurse@hospital.com",
    "first_name": "Jane",
    "last_name": "Smith",
    "hospital_id": 1,
    "is_admin": false
  }
}
```

---

### Update Staff

#### `PUT /api/v1/staff/update/:id`

Update a staff member.

**Authentication:** Bearer Token  
**Tenant Required:** Yes

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | uint | Staff ID |

**Request Body:**
| Field | Type | Required | Validation | Description |
|-------|------|----------|------------|-------------|
| `staff_code` | string | ❌ | max=10 | Staff code |
| `phone_number` | string | ✅ | max=10 | Contact phone |
| `email` | string | ✅ | email format | Email address |
| `first_name` | string | ❌ | max=255 | First name |
| `last_name` | string | ❌ | max=255 | Last name |
| `hospital_id` | uint | ❌ | - | Hospital ID |
| `is_admin` | bool | ❌ | - | Admin privileges |

**Request Example:**
```json
{
  "phone_number": "0898765432",
  "email": "updated@hospital.com",
  "first_name": "Jane",
  "last_name": "Smith-Jones"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "staff updated successfully",
  "data": { ... }
}
```

---

### Delete Staff

#### `DELETE /api/v1/staff/delete/:id`

Delete a staff member. **Admin only.**

**Authentication:** Bearer Token (Admin)  
**Tenant Required:** Yes

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | uint | Staff ID |

**Success Response (200):**
```json
{
  "success": true,
  "message": "staff deleted successfully",
  "data": null
}
```

---

## Patient Endpoints

All patient endpoints require authentication and tenant context.

### Search Patients

#### `GET /api/v1/patient/search`

Search patients by query string.

**Authentication:** Bearer Token  
**Tenant Required:** Yes

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `query` | string | ❌ | Search term (name, national_id, patient_hn, etc.) |

**Request Example:**
```
GET /api/v1/patient/search?query=สมชาย
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "success",
  "data": [
    {
      "ID": 1,
      "CreatedAt": "2024-01-01T00:00:00Z",
      "UpdatedAt": "2024-01-01T00:00:00Z",
      "DeletedAt": null,
      "first_name_th": "สมชาย",
      "last_name_th": "ใจดี",
      "middle_name_th": "",
      "first_name_en": "Somchai",
      "last_name_en": "Jaidee",
      "middle_name_en": "",
      "date_of_birth": "1990-01-15T00:00:00Z",
      "nick_name_th": "ชาย",
      "nick_name_en": "Chai",
      "patient_hn": "BKGH0001-000001",
      "national_id": "1234567890123",
      "passport_id": "AB1234567",
      "phone_number": "0812345678",
      "email": "somchai@email.com",
      "gender": "M",
      "nationality": "Thai",
      "blood_grp": "O"
    }
  ]
}
```

**Empty Result Response (200):**
```json
{
  "success": true,
  "message": "success",
  "data": []
}
```

---

### Create Patient

#### `POST /api/v1/patient/create`

Create a new patient record.

**Authentication:** Bearer Token  
**Tenant Required:** Yes

**Request Body:**
| Field | Type | Required | Validation | Description |
|-------|------|----------|------------|-------------|
| `first_name_th` | string | ✅ | max=255 | Thai first name |
| `last_name_th` | string | ✅ | max=255 | Thai last name |
| `middle_name_th` | string | ❌ | max=255 | Thai middle name |
| `first_name_en` | string | ✅ | max=255 | English first name |
| `last_name_en` | string | ✅ | max=255 | English last name |
| `middle_name_en` | string | ❌ | max=255 | English middle name |
| `date_of_birth` | string | ✅ | YYYY-MM-DD | Date of birth |
| `nick_name_th` | string | ✅ | max=50 | Thai nickname |
| `nick_name_en` | string | ✅ | max=50 | English nickname |
| `national_id` | string | ✅ | max=13 | National ID (13 digits) |
| `passport_id` | string | ✅ | max=13 | Passport number |
| `phone_number` | string | ✅ | max=10 | Phone number |
| `email` | string | ✅ | email format | Email address |
| `gender` | enum | ✅ | M, F, OTHER | Gender |
| `nationality` | string | ✅ | max=100 | Nationality |
| `blood_grp` | enum | ✅ | A, B, O, AB | Blood group |

**Request Example:**
```json
{
  "first_name_th": "สมชาย",
  "last_name_th": "ใจดี",
  "middle_name_th": "",
  "first_name_en": "Somchai",
  "last_name_en": "Jaidee",
  "middle_name_en": "",
  "date_of_birth": "1990-01-15",
  "nick_name_th": "ชาย",
  "nick_name_en": "Chai",
  "national_id": "1234567890123",
  "passport_id": "AB1234567",
  "phone_number": "0812345678",
  "email": "somchai@email.com",
  "gender": "M",
  "nationality": "Thai",
  "blood_grp": "O"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "patient created successfully",
  "data": {
    "ID": 1,
    "CreatedAt": "2024-01-01T00:00:00Z",
    "UpdatedAt": "2024-01-01T00:00:00Z",
    "DeletedAt": null,
    "first_name_th": "สมชาย",
    "last_name_th": "ใจดี",
    "patient_hn": "BKGH0001-000001",
    ...
  }
}
```

---

### Update Patient (Full)

#### `PUT /api/v1/patient/update/:id`

Full update of a patient record. All required fields must be provided.

**Authentication:** Bearer Token  
**Tenant Required:** Yes

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | uint | Patient ID |

**Request Body:** Same as Create Patient

**Success Response (200):**
```json
{
  "success": true,
  "message": "patient updated successfully",
  "data": { ... }
}
```

**Error Responses:**
| Status | Error |
|--------|-------|
| 400 | `invalid id` |
| 404 | `patient not found` |
| 409 | `duplicate patient entry` |

---

### Update Patient (Partial)

#### `PATCH /api/v1/patient/update/:id`

Partial update of a patient record. Only provided fields will be updated.

**Authentication:** Bearer Token  
**Tenant Required:** Yes

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | uint | Patient ID |

**Request Body:**
All fields are optional. Only include fields you want to update.

| Field | Type | Validation | Description |
|-------|------|------------|-------------|
| `first_name_th` | string | max=255 | Thai first name |
| `last_name_th` | string | max=255 | Thai last name |
| `middle_name_th` | string | max=255 | Thai middle name |
| `first_name_en` | string | max=255 | English first name |
| `last_name_en` | string | max=255 | English last name |
| `middle_name_en` | string | max=255 | English middle name |
| `date_of_birth` | string | YYYY-MM-DD | Date of birth |
| `nick_name_th` | string | max=50 | Thai nickname |
| `nick_name_en` | string | max=50 | English nickname |
| `national_id` | string | max=13 | National ID |
| `passport_id` | string | max=13 | Passport number |
| `phone_number` | string | max=10 | Phone number |
| `email` | string | email format | Email address |
| `gender` | enum | M, F, OTHER | Gender |
| `nationality` | string | max=100 | Nationality |
| `blood_grp` | enum | A, B, O, AB | Blood group |

**Request Example:**
```json
{
  "phone_number": "0898765432",
  "email": "newemail@email.com"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "patient updated successfully",
  "data": { ... }
}
```

---

### Delete Patient

#### `DELETE /api/v1/patient/delete/:id`

Delete a patient record (soft delete).

**Authentication:** Bearer Token  
**Tenant Required:** Yes

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | uint | Patient ID |

**Success Response (200):**
```json
{
  "success": true,
  "message": "patient deleted successfully",
  "data": null
}
```

**Error Response (404):**
```json
{
  "success": false,
  "error": "patient not found"
}
```

---

## Error Codes

| HTTP Status | Error Message | Description |
|-------------|---------------|-------------|
| 400 | `invalid id` | Invalid ID parameter |
| 400 | `invalid input` | Request validation failed |
| 400 | `invalid date format` | Date not in YYYY-MM-DD format |
| 400 | `date of birth cannot be in the future` | Future date provided |
| 400 | `date of birth is too old` | Date exceeds 150 years |
| 400 | `tenant context required` | Missing tenant information |
| 400 | `invalid tenant schema` | Invalid schema name |
| 401 | `unauthorized` | Missing or invalid token |
| 401 | `invalid credentials` | Wrong username/password |
| 403 | `forbidden` | Insufficient permissions |
| 404 | `not found` | Resource not found |
| 409 | `duplicate entry` | Unique constraint violation |
| 500 | `internal server error` | Server error |

---

## Data Types

### Gender Enum
| Value | Description |
|-------|-------------|
| `M` | Male |
| `F` | Female |
| `OTHER` | Other |

### Blood Group Enum
| Value | Description |
|-------|-------------|
| `A` | Blood type A |
| `B` | Blood type B |
| `O` | Blood type O |
| `AB` | Blood type AB |

---

## Authentication

### JWT Token

Include the JWT token in the `Authorization` header:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Token Claims

| Claim | Type | Description |
|-------|------|-------------|
| `user_id` | uint | Staff ID |
| `username` | string | Staff username |
| `is_admin` | bool | Admin privileges |
| `hospital_id` | uint | Hospital ID |
| `exp` | int64 | Expiration timestamp |

---

## Rate Limiting

Currently not implemented.

---

## Examples

### cURL Examples

**Login:**
```bash
curl -X POST "https://bangkok.his.com/api/v1/staff/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123",
    "hospital_id": 1
  }'
```

**Get All Staff (Authenticated):**
```bash
curl -X GET "https://bangkok.his.com/api/v1/staff/" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Create Patient:**
```bash
curl -X POST "https://bangkok.his.com/api/v1/patient/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -d '{
    "first_name_th": "สมชาย",
    "last_name_th": "ใจดี",
    "first_name_en": "Somchai",
    "last_name_en": "Jaidee",
    "date_of_birth": "1990-01-15",
    "nick_name_th": "ชาย",
    "nick_name_en": "Chai",
    "national_id": "1234567890123",
    "passport_id": "AB1234567",
    "phone_number": "0812345678",
    "email": "somchai@email.com",
    "gender": "M",
    "nationality": "Thai",
    "blood_grp": "O"
  }'
```

**Search Patients:**
```bash
curl -X GET "https://bangkok.his.com/api/v1/patient/search?query=สมชาย" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Partial Update Patient:**
```bash
curl -X PATCH "https://bangkok.his.com/api/v1/patient/update/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -d '{
    "phone_number": "0898765432"
  }'
```

---

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2024-01-01 | Initial API release |
