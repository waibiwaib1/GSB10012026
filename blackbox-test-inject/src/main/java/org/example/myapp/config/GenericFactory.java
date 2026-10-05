package org.example.myapp.config;

import io.avaje.inject.Bean;
import io.avaje.inject.Factory;

@Factory
public class GenericFactory {

  @Bean
  public ParamType<MyDoc> buildParam() {
    return new ParamType<>();
  }

  public static class ParamType<T> {}

  public static class MyDoc {}
}
