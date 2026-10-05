package org.example.myapp.config;

import jakarta.inject.Singleton;

@Singleton
public class ParamConsumer {

  public final GenericFactory.ParamType<GenericFactory.MyDoc> param;

  public ParamConsumer(GenericFactory.ParamType<GenericFactory.MyDoc> param) {
    this.param = param;
  }
}
