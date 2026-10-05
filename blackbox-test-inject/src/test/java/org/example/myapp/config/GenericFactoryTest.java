package org.example.myapp.config;

import io.avaje.inject.BeanScope;
import io.avaje.inject.spi.GenericType;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertNotNull;

class GenericFactoryTest {

  @Test
  void parameterizedTypeLookup() {
    try (BeanScope scope = BeanScope.builder().build()) {
      var bean = scope.get(new GenericType<GenericFactory.ParamType<GenericFactory.MyDoc>>() {}, null);
      assertNotNull(bean);
    }
  }

  @Test
  void parameterizedTypeInjection() {
    try (BeanScope scope = BeanScope.builder().build()) {
      var consumer = scope.get(ParamConsumer.class);
      assertNotNull(consumer.param);
    }
  }
}
