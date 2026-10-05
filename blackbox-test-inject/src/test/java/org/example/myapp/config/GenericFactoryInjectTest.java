package org.example.myapp.config;

import io.avaje.inject.test.InjectTest;
import jakarta.inject.Inject;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertNotNull;

@InjectTest
class GenericFactoryInjectTest {

  @Inject GenericFactory.ParamType<GenericFactory.MyDoc> param;

  @Test
  void injected() {
    assertNotNull(param);
  }
}
