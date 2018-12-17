# Data Types

- `Hash` : SHA256 hash
- `UUID` : universally unique identifier formatted similar to `a47178c1-3367-440d-b35e-a2a0a644a586`
- `URL` : resource location

# Admin

## Ping

Used as a healthcheck.

**URL** : `/ping/`

**Method** : `GET`

### Success Response

**Code** : `200 OK`

**Content example**

```
Pong
```

# Document

## Data Types

- `DocumentState` : enum of strings, ['pending', 'batched', 'failed']

## Get Document

Get an existing document.

**URL** : `/document/{UUID:id}`

- id : UUID of the document

**Method** : `GET`

### Success Response

**Code** : `200 OK`

**Data Format**

```
{
    "hash": Hash,
    "id": UUID,
    "stamp": (URL | Null),
    "state": DocumentState
}
```

**Content example**

```json
{
    "hash": "f3ac645b600b5dfef586fc8f5f1fae27b9f67c90",
    "id": "febc23bd-3474-4f2e-85eb-ba109ce62ed5",
    "stamp": "localhost:9000/stamp/3088f26a-c1a3-4fad-8f73-9b1b3d5a0dbf",
    "state": "batched"
}
```

## Create Document

Create a new document.

**URL** : `/document`

**Method** : `POST`

**Data Format**


```txt
{
    "hash": Hash
}
```

**Data example**

```json
{
    "hash": "0d6772a577eb705cfce697b9e5da7a4034944bb3"
}
```

### Success Response

**Code** : `200 OK`

**Content example**

```txt
e64fb3ad-cc2e-40ea-ad28-1b67d8b741ea
```

# Stamp

## Data Types

- `MerkleTree` :
```
{
    "root": MerkleNode
}
```
- `MerkleNode` :
```
{
    "hash": Hash,
    "left": (MerkleNode | Null),
    "right": (MerkleNode | Null)
}
```
- `StampState` : enum of strings ['pending', 'processing', 'sent', 'confirmed', 'failed']

## Get Stamp

Get an existing stamp.

**URL** : `/stamp/{id}`

- id : UUID of the document

**Method** : `GET`

## Success Response

**Code** : `200 OK`

**Content Format**
```
{
    "id": UUID,
    "merkleTree": MerkleTree,
    "state": StampState,
    "txhash": (TransactionHash | Null)
}
```

**Content example**

```
{
    "id": "a47178c1-3367-440d-b35e-a2a0a644a586",
    "merkleTree": {
        "root": {
            "hash": "1aed2ae581e8435bc7f00c76085620509bf656953f57170b09e0c84810ba20a9",
            "left": {
                "hash": "1d6772a577eb705cfce697b9e5da7a4034944bb3",
                "left": null,
                "right": null
            },
            "right": null
        }
    },
    "state": "sent",
    "txhash": "0xcbb8965477ed40a316ccb101e8a01697d55b7dbb1a4f68992514144a8f9e6255"
}
```
