# Configuration

Displayed all configuration options.

----------

## First level

All keys of first level configuration.

* `debug` - bool
* `indexes` - IndexesStruct

----------

## IndexesStruct

Configuration of IndexesStruct.

* `user` - string - Default: `users`
  * Sets the index name to the users in ES.

```yaml
indexes:
  user: users
```

