import {Attribute, Entity, INDEX_TYPE} from '@typedorm/common';
import {table} from './table';

export interface UserSparseIndexPrimaryKey {
  id: string;
}

@Entity({
  table,
  name: 'user-sparse-index',
  primaryKey: {
    partitionKey: 'USER#{{id}}',
    sortKey: 'USER#{{id}}',
  },
  indexes: {
    GSI2: {
      partitionKey: 'USER#STATUS#{{status}}',
      sortKey: 'USER#PUB_ID#{{pubId}}',
      type: INDEX_TYPE.GSI,
      isSparse: true,
    },
  },
})
export class UserSparseIndex implements UserSparseIndexPrimaryKey {
  @Attribute()
  id: string;

  @Attribute()
  name: string;

  @Attribute()
  status: string;

  @Attribute()
  pubId?: string;
}
